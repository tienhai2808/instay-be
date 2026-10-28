package supportticket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/data"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/mapper"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/internal/application/usecase/branchscope"
	"github.com/InstaySystem/is_v2-be/internal/domain/constants"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
	"github.com/InstaySystem/is_v2-be/pkg/utils"
)

type supportTicketUseCaseImpl struct {
	logPro              port.LogProvider
	idPro               port.IDProvider
	txPro               port.TxProvider
	mqPro               port.MessageQueueProvider
	supportTicketRepo   repository.SupportTicketRepository
	bookingRoomRepo     repository.BookingRoomRepository
	supportCategoryRepo repository.SupportCategoryRepository
	notificationRepo    repository.NotificationRepository
	userRepo            repository.UserRepository
	deviceTokenRepo     repository.DeviceTokenRepository
}

func NewSupportTicketUseCase(
	logPro port.LogProvider,
	idPro port.IDProvider,
	txPro port.TxProvider,
	mqPro port.MessageQueueProvider,
	supportTicketRepo repository.SupportTicketRepository,
	bookingRoomRepo repository.BookingRoomRepository,
	supportCategoryRepo repository.SupportCategoryRepository,
	notificationRepo repository.NotificationRepository,
	userRepo repository.UserRepository,
	deviceTokenRepo repository.DeviceTokenRepository,
) SupportTicketUseCase {
	return &supportTicketUseCaseImpl{
		logPro,
		idPro,
		txPro,
		mqPro,
		supportTicketRepo,
		bookingRoomRepo,
		supportCategoryRepo,
		notificationRepo,
		userRepo,
		deviceTokenRepo,
	}
}

func (u *supportTicketUseCaseImpl) Create(ctx context.Context, schema string, bookingRoomID int64, req request.CreateSupportTicketRequest) (int64, error) {
	bookingRoom, err := u.bookingRoomRepo.FindByIDWithRoom(ctx, schema, bookingRoomID)
	if err != nil {
		u.logPro.Error("Find booking room by id failed", port.String("schema", schema), port.Int64("booking_room_id", bookingRoomID), port.Error(err))
		return 0, err
	}
	if bookingRoom == nil {
		return 0, exceptions.ErrUnauth
	}

	supCat, err := u.supportCategoryRepo.FindByID(ctx, schema, req.SupportCategoryID)
	if err != nil {
		u.logPro.Error("Find support category by id failed", port.String("schema", schema), port.Int64("support_category_id", req.SupportCategoryID), port.Error(err))
		return 0, err
	}
	if supCat == nil {
		return 0, exceptions.ErrSupportCategoryNotFound
	}

	id, err := u.idPro.NextID()
	if err != nil {
		u.logPro.Error("Generate order service id failed", port.String("schema", schema), port.Error(err))
		return 0, err
	}

	supTicket := model.NewSupportTicket(id, bookingRoomID, req.Content, model.SupportTicketStatusPending, req.Priority, req.SupportCategoryID)

	if err = u.txPro.WithinTx(ctx, func(ctx context.Context) error {
		if err = u.supportTicketRepo.Create(ctx, schema, supTicket); err != nil {
			u.logPro.Error("Create support ticket failed", port.String("schema", schema), port.Error(err))
			return err
		}

		nid, err := u.idPro.NextID()
		if err != nil {
			u.logPro.Error("Generate notification id failed", port.String("schema", schema), port.Error(err))
			return err
		}

		notificationContentVI := fmt.Sprintf("Phòng %s yêu cầu hỗ trợ về %s", bookingRoom.Room.Name, supCat.Name)
		notificationContentEN := fmt.Sprintf("Room %s requires assistance with %s", bookingRoom.Room.Name, supCat.Name)

		notification := model.NewNotification(nid, &supCat.DepartmentID, &bookingRoomID, model.NotificationReceiverTypeStaff, model.NotificationContentTypeSupportTicket, strconv.FormatInt(id, 10), notificationContentVI, notificationContentEN, false)

		if err = u.notificationRepo.Create(ctx, schema, notification); err != nil {
			u.logPro.Error("Create notification failed", port.String("schema", schema), port.Error(err))
			return err
		}

		users, err := u.userRepo.FindAllActiveByDepartmentID(ctx, schema, supCat.DepartmentID)
		if err != nil {
			u.logPro.Error("Find all active user by department id failed", port.String("schema", schema), port.Error(err))
			return err
		}

		userIDs := userIDsForBranch(users, bookingRoom.Room.Floor.BranchID)

		deviceTokens, err := u.deviceTokenRepo.FindAllByUserIDs(ctx, schema, userIDs)
		if err != nil {
			u.logPro.Error("Find all device token by user ids failed", port.String("schema", schema), port.Error(err))
			return err
		}

		tokens := make([]string, 0, len(deviceTokens))
		for _, deviceToken := range deviceTokens {
			tokens = append(tokens, deviceToken.Token)
		}

		pnData := map[string]string{
			"code":         strconv.Itoa(constants.CodePushNotificationSupportTicketNew),
			"content_key":  strconv.FormatInt(id, 10),
			"content_type": string(model.NotificationContentTypeSupportTicket),
		}

		u.pushNotification(schema, tokens, pnData)

		return nil
	}); err != nil {
		return 0, err
	}

	return id, nil
}

func (u *supportTicketUseCaseImpl) ListForGuest(ctx context.Context, schema string, bookingRoomID int64) ([]*response.SimpleSupportTicketResponse, error) {
	supTickets, err := u.supportTicketRepo.FindAllByBookingRoomIDWithSupportCategory(ctx, schema, bookingRoomID)
	if err != nil {
		u.logPro.Error("Find all support ticket by booking room id failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}

	return mapper.ToSimpleSupportTicketsResponse(supTickets), nil
}

func (u *supportTicketUseCaseImpl) UpdateStatusForGuest(ctx context.Context, schema string, supportTicketID, bookingRoomID int64, req request.UpdateSupportTicketStatusForGuestRequest) error {
	bookingRoom, err := u.bookingRoomRepo.FindByIDWithRoom(ctx, schema, bookingRoomID)
	if err != nil {
		u.logPro.Error("Find booking room by id failed", port.String("schema", schema), port.Int64("booking_room_id", bookingRoomID), port.Error(err))
		return err
	}
	if bookingRoom == nil {
		return exceptions.ErrUnauth
	}

	supportTicket, err := u.supportTicketRepo.FindByIDWithSupportCategory(ctx, schema, supportTicketID)
	if err != nil {
		u.logPro.Error("Find support ticket by id failed", port.String("schema", schema), port.Int64("support_ticket_id", supportTicketID), port.Error(err))
		return err
	}
	if supportTicket == nil {
		return exceptions.ErrSupportTicketNotFound
	}

	if supportTicket.BookingRoomID != bookingRoomID {
		return exceptions.ErrForbidden
	}

	if supportTicket.Status != model.SupportTicketStatusPending || req.Status != model.SupportTicketStatusCancelled {
		return exceptions.ErrInvalidStatus
	}

	if err = u.txPro.WithinTx(ctx, func(ctx context.Context) error {
		if err := u.supportTicketRepo.Update(ctx, schema, supportTicketID, map[string]any{"status": model.SupportTicketStatusCancelled}); err != nil {
			if errors.Is(err, exceptions.ErrSupportTicketNotFound) {
				return err
			}
			u.logPro.Error("Update support ticket failed", port.String("schema", schema), port.Int64("support_ticket_id", supportTicketID), port.Error(err))
			return err
		}

		deptID := supportTicket.SupportCategory.DepartmentID

		id, err := u.idPro.NextID()
		if err != nil {
			u.logPro.Error("Generate notification id failed", port.String("schema", schema), port.Error(err))
			return err
		}

		notificationContentVI := fmt.Sprintf("Phòng %s đã hủy phiếu hỗ trợ về %s", bookingRoom.Room.Name, supportTicket.SupportCategory.Name)
		notificationContentEN := fmt.Sprintf("Room %s cancelled support ticket about %s", bookingRoom.Room.Name, supportTicket.SupportCategory.Name)

		notification := model.NewNotification(id, &deptID, &bookingRoomID, model.NotificationReceiverTypeStaff, model.NotificationContentTypeSupportTicket, strconv.FormatInt(supportTicketID, 10), notificationContentVI, notificationContentEN, false)

		if err = u.notificationRepo.Create(ctx, schema, notification); err != nil {
			u.logPro.Error("Create notification failed", port.String("schema", schema), port.Error(err))
			return err
		}

		users, err := u.userRepo.FindAllActiveByDepartmentID(ctx, schema, deptID)
		if err != nil {
			u.logPro.Error("Find all active user by department id failed", port.String("schema", schema), port.Error(err))
			return err
		}

		userIDs := userIDsForBranch(users, bookingRoom.Room.Floor.BranchID)

		deviceTokens, err := u.deviceTokenRepo.FindAllByUserIDs(ctx, schema, userIDs)
		if err != nil {
			u.logPro.Error("Find all device token by user ids failed", port.String("schema", schema), port.Error(err))
			return err
		}

		tokens := make([]string, 0, len(deviceTokens))
		for _, deviceToken := range deviceTokens {
			tokens = append(tokens, deviceToken.Token)
		}

		pnData := map[string]string{
			"code":         strconv.Itoa(constants.CodePushNotificationSupportTicketCancelled),
			"content_key":  strconv.FormatInt(supportTicketID, 10),
			"content_type": string(model.NotificationContentTypeSupportTicket),
		}

		u.pushNotification(schema, tokens, pnData)

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (u *supportTicketUseCaseImpl) ListForAdmin(ctx context.Context, schema string, userID int64, departmentID *int64, query request.SupportTicketPaginationQuery) ([]*response.SupportTicketResponse, *response.MetaResponse, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	supportTicketQuery := mapper.ToSupportTicketQuery(query)
	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return nil, nil, err
	}
	if !scope.IsGlobal() {
		supportTicketQuery.BranchIDs = scope.BranchIDs
	}

	supTickets, total, err := u.supportTicketRepo.FindAllWithSupportCategoryPaginated(ctx, schema, departmentID, supportTicketQuery)
	if err != nil {
		u.logPro.Error("Find all support ticket paginated failed", port.String("schema", schema), port.Error(err))
		return nil, nil, err
	}

	meta := mapper.ToMetaResponse(total, query.Page, query.Limit)

	return mapper.ToSupportTicketsResponse(supTickets), meta, nil
}

func (u *supportTicketUseCaseImpl) UpdateStatusForAdmin(ctx context.Context, schema string, supportTicketID, userID int64, req request.UpdateSupportTicketStatusForAdminRequest) error {
	supportTicket, err := u.supportTicketRepo.FindByIDWithSupportCategory(ctx, schema, supportTicketID)
	if err != nil {
		u.logPro.Error("Find support ticket by id failed", port.String("schema", schema), port.Int64("support_ticket_id", supportTicketID), port.Error(err))
		return err
	}
	if supportTicket == nil {
		return exceptions.ErrSupportTicketNotFound
	}
	bookingRoom, err := u.bookingRoomRepo.FindByIDWithRoom(ctx, schema, supportTicket.BookingRoomID)
	if err != nil {
		u.logPro.Error("Find booking room by id failed", port.String("schema", schema), port.Int64("booking_room_id", supportTicket.BookingRoomID), port.Error(err))
		return err
	}
	if err := u.ensureBookingRoomInUserScope(ctx, schema, userID, bookingRoom); err != nil {
		return err
	}

	switch supportTicket.Status {
	case model.SupportTicketStatusPending:
		if req.Status != model.SupportTicketStatusAccepted {
			return exceptions.ErrInvalidStatus
		}
	case model.SupportTicketStatusAccepted:
		if req.Status != model.SupportTicketStatusDone {
			return exceptions.ErrInvalidStatus
		}
	default:
		return exceptions.ErrInvalidStatus
	}

	if err = u.txPro.WithinTx(ctx, func(ctx context.Context) error {
		updateData := map[string]any{
			"status":        req.Status,
			"updated_by_id": userID,
		}
		if req.Status == model.SupportTicketStatusAccepted {
			updateData["resolution_seconds"] = nil
		}
		if req.Status == model.SupportTicketStatusDone {
			durationSeconds := max(int64(time.Since(supportTicket.UpdatedAt).Seconds()), 0)
			updateData["resolution_seconds"] = durationSeconds
		}

		if err := u.supportTicketRepo.Update(ctx, schema, supportTicketID, updateData); err != nil {
			if errors.Is(err, exceptions.ErrSupportTicketNotFound) {
				return err
			}
			u.logPro.Error("Update support ticket failed", port.String("schema", schema), port.Int64("support_ticket_id", supportTicketID), port.Error(err))
			return err
		}

		id, err := u.idPro.NextID()
		if err != nil {
			u.logPro.Error("Generate notification id failed", port.String("schema", schema), port.Error(err))
			return err
		}

		categoryName := supportTicket.SupportCategory.Name

		notificationContentVI := fmt.Sprintf("Phiếu hỗ trợ về %s đã được xác nhận", categoryName)
		notificationContentEN := fmt.Sprintf("Support ticket about %s has been accepted", categoryName)

		notification := model.NewNotification(id, nil, &supportTicket.BookingRoomID, model.NotificationReceiverTypeGuest, model.NotificationContentTypeSupportTicket, strconv.FormatInt(supportTicketID, 10), notificationContentVI, notificationContentEN, false)

		if err = u.notificationRepo.Create(ctx, schema, notification); err != nil {
			u.logPro.Error("Create notification failed", port.String("schema", schema), port.Error(err))
			return err
		}

		deviceTokens, err := u.deviceTokenRepo.FindAllByBookingRoomID(ctx, schema, supportTicket.BookingRoomID)
		if err != nil {
			u.logPro.Error("Find all device token by booking room id failed", port.String("schema", schema), port.Error(err))
			return err
		}

		tokens := make([]string, 0, len(deviceTokens))
		for _, deviceToken := range deviceTokens {
			tokens = append(tokens, deviceToken.Token)
		}

		pnData := map[string]string{
			"code":         strconv.Itoa(constants.CodePushNotificationSupportTicketConfirmed),
			"content_key":  strconv.FormatInt(supportTicketID, 10),
			"content_type": string(model.NotificationContentTypeSupportTicket),
		}

		u.pushNotification(schema, tokens, pnData)

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (u *supportTicketUseCaseImpl) pushNotification(schema string, tokens []string, pnData map[string]string) {
	if len(tokens) == 0 {
		return
	}

	pushNotificationMsg := data.PushNotificationMessage{
		Schema: schema,
		Tokens: tokens,
		Data:   pnData,
	}

	go func(msg data.PushNotificationMessage) {
		body, err := json.Marshal(msg)
		if err != nil {
			u.logPro.Error("Json marshal push notification message failed", port.String("schema", schema), port.Error(err))
		}

		routingKey := utils.GetRoutingKey(schema, port.RoutingKeyPushNotification)
		if err = u.mqPro.PublishMessage(port.ExchangePushNotification, routingKey, body); err != nil {
			u.logPro.Error("Publish push notification message failed", port.String("schema", schema), port.Error(err))
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		if err = u.deviceTokenRepo.Updates(ctx, schema, tokens, map[string]any{"last_used_at": time.Now()}); err != nil {
			u.logPro.Error("Update device token by tokens failed", port.String("schema", schema), port.Error(err))
		}
	}(pushNotificationMsg)
}

func (u *supportTicketUseCaseImpl) ensureBookingRoomInUserScope(ctx context.Context, schema string, userID int64, bookingRoom *model.BookingRoom) error {
	if bookingRoom == nil || bookingRoom.Room == nil || bookingRoom.Room.Floor == nil {
		return exceptions.ErrBookingRoomNotFound
	}

	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return err
	}

	return scope.Ensure(bookingRoom.Room.Floor.BranchID)
}

func userIDsForBranch(users []*model.User, branchID int64) []int64 {
	userIDs := make([]int64, 0, len(users))
	for _, user := range users {
		if user == nil {
			continue
		}
		if len(user.Branches) == 0 {
			userIDs = append(userIDs, user.ID)
			continue
		}
		for _, branch := range user.Branches {
			if branch != nil && branch.ID == branchID {
				userIDs = append(userIDs, user.ID)
				break
			}
		}
	}

	return userIDs
}
