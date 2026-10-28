package orderservice

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

type orderServiceUseCaseImpl struct {
	logPro           port.LogProvider
	idPro            port.IDProvider
	txPro            port.TxProvider
	pnPro            port.PushNotificationProvider
	mqPro            port.MessageQueueProvider
	serviceRepo      repository.ServiceRepository
	bookingRoomRepo  repository.BookingRoomRepository
	orderServiceRepo repository.OrderServiceRepository
	notificationRepo repository.NotificationRepository
	userRepo         repository.UserRepository
	deviceTokenRepo  repository.DeviceTokenRepository
}

func NewOrderServiceUseCase(
	log port.LogProvider,
	idGen port.IDProvider,
	txPro port.TxProvider,
	pnPro port.PushNotificationProvider,
	mqPro port.MessageQueueProvider,
	serviceRepo repository.ServiceRepository,
	bookingRoomRepo repository.BookingRoomRepository,
	orderServiceRepo repository.OrderServiceRepository,
	notificationRepo repository.NotificationRepository,
	userRepo repository.UserRepository,
	deviceTokenRepo repository.DeviceTokenRepository,
) OrderServiceUseCase {
	return &orderServiceUseCaseImpl{
		log,
		idGen,
		txPro,
		pnPro,
		mqPro,
		serviceRepo,
		bookingRoomRepo,
		orderServiceRepo,
		notificationRepo,
		userRepo,
		deviceTokenRepo,
	}
}

func (u *orderServiceUseCaseImpl) CreateForGuest(ctx context.Context, schema string, bookingRoomID int64, req request.CreateOrderServiceForGuestRequest) (int64, error) {
	bookingRoom, err := u.bookingRoomRepo.FindByIDWithRoom(ctx, schema, bookingRoomID)
	if err != nil {
		u.logPro.Error("Find booking room by id failed", port.String("schema", schema), port.Int64("booking_room_id", bookingRoomID), port.Error(err))
		return 0, err
	}
	if bookingRoom == nil {
		return 0, exceptions.ErrUnauth
	}

	service, err := u.serviceRepo.FindByIDWithServiceType(ctx, schema, req.ServiceID)
	if err != nil {
		u.logPro.Error("Find service by id failed", port.String("schema", schema), port.Int64("service_id", req.ServiceID), port.Error(err))
		return 0, err
	}
	if service == nil {
		return 0, exceptions.ErrServiceNotFound
	}
	if err = u.ensureServiceInBookingRoomBranch(ctx, schema, req.ServiceID, bookingRoom); err != nil {
		return 0, err
	}

	id, err := u.idPro.NextID()
	if err != nil {
		u.logPro.Error("Generate order service id failed", port.String("schema", schema), port.Error(err))
		return 0, err
	}

	unitPrice := service.UnitPriceAt(req.EstimatedAt)
	totalPrice := float64(req.Quantity) * unitPrice

	orderService := model.NewOrderService(id, req.ServiceID, bookingRoomID, req.Quantity, unitPrice, totalPrice, model.OrderServiceStatusPending, req.GuestNote, nil, nil, req.EstimatedAt, nil, nil)

	if err = u.txPro.WithinTx(ctx, func(ctx context.Context) error {
		if err = u.orderServiceRepo.Create(ctx, schema, orderService); err != nil {
			u.logPro.Error("Create order service failed", port.String("schema", schema), port.Error(err))
			return err
		}

		nid, err := u.idPro.NextID()
		if err != nil {
			u.logPro.Error("Generate notification id failed", port.String("schema", schema), port.Error(err))
			return err
		}

		notificationContentVI := fmt.Sprintf("Phòng %s đã đặt %d %s", bookingRoom.Room.Name, req.Quantity, service.Name)
		notificationContentEN := fmt.Sprintf("Room %s ordered %d %s", bookingRoom.Room.Name, req.Quantity, service.Name)

		notification := model.NewNotification(nid, &service.ServiceType.DepartmentID, &bookingRoomID, model.NotificationReceiverTypeStaff, model.NotificationContentTypeOrderService, strconv.FormatInt(id, 10), notificationContentVI, notificationContentEN, false)

		if err = u.notificationRepo.Create(ctx, schema, notification); err != nil {
			u.logPro.Error("Create notification failed", port.String("schema", schema), port.Error(err))
			return err
		}

		users, err := u.userRepo.FindAllActiveByDepartmentID(ctx, schema, service.ServiceType.DepartmentID)
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
			"code":         strconv.Itoa(constants.CodePushNotificationOrderServiceNew),
			"content_key":  strconv.FormatInt(id, 10),
			"content_type": string(model.NotificationContentTypeOrderService),
		}

		u.pushNotification(schema, tokens, pnData)

		return nil
	}); err != nil {
		return 0, err
	}

	return id, nil
}

func (u *orderServiceUseCaseImpl) CreateForAdmin(ctx context.Context, schema string, userID int64, req request.CreateOrderServiceForAdminRequest) (int64, error) {
	bookingRoom, err := u.bookingRoomRepo.FindByIDWithRoom(ctx, schema, req.BookingRoomID)
	if err != nil {
		u.logPro.Error("Find booking room by id failed", port.String("schema", schema), port.Int64("booking_room_id", req.BookingRoomID), port.Error(err))
		return 0, err
	}
	if bookingRoom == nil {
		return 0, exceptions.ErrUnauth
	}
	if err := u.ensureBookingRoomInUserScope(ctx, schema, userID, bookingRoom); err != nil {
		return 0, err
	}

	service, err := u.serviceRepo.FindByIDWithServiceType(ctx, schema, req.ServiceID)
	if err != nil {
		u.logPro.Error("Find service by id failed", port.String("schema", schema), port.Int64("service_id", req.ServiceID), port.Error(err))
		return 0, err
	}
	if service == nil {
		return 0, exceptions.ErrServiceNotFound
	}
	if err = u.ensureServiceInBookingRoomBranch(ctx, schema, req.ServiceID, bookingRoom); err != nil {
		return 0, err
	}

	id, err := u.idPro.NextID()
	if err != nil {
		u.logPro.Error("Generate order service id failed", port.String("schema", schema), port.Error(err))
		return 0, err
	}

	unitPrice := service.UnitPriceAt(req.EstimatedAt)
	totalPrice := float64(req.Quantity) * unitPrice

	orderService := model.NewOrderService(id, req.ServiceID, req.BookingRoomID, req.Quantity, unitPrice, totalPrice, model.OrderServiceStatusPending, req.GuestNote, nil, nil, req.EstimatedAt, &userID, &userID)

	if err = u.txPro.WithinTx(ctx, func(ctx context.Context) error {
		if err = u.orderServiceRepo.Create(ctx, schema, orderService); err != nil {
			u.logPro.Error("Create order service failed", port.String("schema", schema), port.Error(err))
			return err
		}

		nid, err := u.idPro.NextID()
		if err != nil {
			u.logPro.Error("Generate notification id failed", port.String("schema", schema), port.Error(err))
			return err
		}

		notificationContentVI := fmt.Sprintf("Phòng %s đã đặt %d %s", bookingRoom.Room.Name, req.Quantity, service.Name)
		notificationContentEN := fmt.Sprintf("Room %s ordered %d %s", bookingRoom.Room.Name, req.Quantity, service.Name)

		notification := model.NewNotification(nid, &service.ServiceType.DepartmentID, &req.BookingRoomID, model.NotificationReceiverTypeStaff, model.NotificationContentTypeOrderService, strconv.FormatInt(id, 10), notificationContentVI, notificationContentEN, false)

		if err = u.notificationRepo.Create(ctx, schema, notification); err != nil {
			u.logPro.Error("Create notification failed", port.String("schema", schema), port.Error(err))
			return err
		}

		users, err := u.userRepo.FindAllActiveByDepartmentID(ctx, schema, service.ServiceType.DepartmentID)
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
			"code":         strconv.Itoa(constants.CodePushNotificationOrderServiceNew),
			"content_key":  strconv.FormatInt(id, 10),
			"content_type": string(model.NotificationContentTypeOrderService),
		}

		u.pushNotification(schema, tokens, pnData)

		return nil
	}); err != nil {
		return 0, err
	}

	return id, nil
}

func (u *orderServiceUseCaseImpl) GetByIDForAdmin(ctx context.Context, schema string, orderServiceID, userID int64) (*response.OrderServiceResponse, error) {
	orderService, err := u.orderServiceRepo.FindByIDWithDetails(ctx, schema, orderServiceID)
	if err != nil {
		u.logPro.Error("Find order service by id failed", port.String("schema", schema), port.Int64("order_service_id", orderServiceID), port.Error(err))
		return nil, err
	}
	if orderService == nil {
		return nil, exceptions.ErrOrderServiceNotFound
	}
	if err := u.ensureBookingRoomInUserScope(ctx, schema, userID, orderService.BookingRoom); err != nil {
		return nil, err
	}

	return mapper.ToOrderServiceResponse(orderService), nil
}

func (u *orderServiceUseCaseImpl) ListForGuest(ctx context.Context, schema string, bookingRoomID int64) ([]*response.SimpleOrderServiceResponse, error) {
	orderServices, err := u.orderServiceRepo.FindAllByBookingRoomIDWithServiceDetails(ctx, schema, bookingRoomID)
	if err != nil {
		u.logPro.Error("Find all order service by booking room id failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}

	return mapper.ToSimpleOrderServicesResponse(orderServices), nil
}

func (u *orderServiceUseCaseImpl) UpdateStatusForGuest(ctx context.Context, schema string, orderServiceID, bookingRoomID int64, req request.UpdateOrderServiceStatusForGuestRequest) error {
	bookingRoom, err := u.bookingRoomRepo.FindByIDWithRoom(ctx, schema, bookingRoomID)
	if err != nil {
		u.logPro.Error("Find booking room by id failed", port.String("schema", schema), port.Int64("booking_room_id", bookingRoomID), port.Error(err))
		return err
	}
	if bookingRoom == nil {
		return exceptions.ErrUnauth
	}

	orderService, err := u.orderServiceRepo.FindByIDWithServiceDetails(ctx, schema, orderServiceID)
	if err != nil {
		u.logPro.Error("Find order service by id failed", port.String("schema", schema), port.Int64("order_service_id", orderServiceID), port.Error(err))
		return err
	}
	if orderService == nil {
		return exceptions.ErrOrderServiceNotFound
	}

	if orderService.BookingRoomID != bookingRoomID {
		return exceptions.ErrForbidden
	}

	if orderService.Status != model.OrderServiceStatusPending || req.Status != model.OrderServiceStatusCancelled {
		return exceptions.ErrInvalidStatus
	}

	if err = u.txPro.WithinTx(ctx, func(ctx context.Context) error {
		updateData := map[string]any{
			"status": model.OrderServiceStatusCancelled,
		}
		if req.CancelReason != nil {
			updateData["cancel_reason"] = req.CancelReason
		}

		if err := u.orderServiceRepo.Update(ctx, schema, orderServiceID, updateData); err != nil {
			if errors.Is(err, exceptions.ErrOrderServiceNotFound) {
				return err
			}
			u.logPro.Error("Update order service failed", port.String("schema", schema), port.Int64("order_service_id", orderServiceID), port.Error(err))
			return err
		}

		deptID := orderService.Service.ServiceType.DepartmentID

		id, err := u.idPro.NextID()
		if err != nil {
			u.logPro.Error("Generate notification id failed", port.String("schema", schema), port.Error(err))
			return err
		}

		notificationContentVI := fmt.Sprintf("Phòng %s đã hủy %d %s", bookingRoom.Room.Name, orderService.Quantity, orderService.Service.Name)
		notificationContentEN := fmt.Sprintf("Room %s cancelled %d %s", bookingRoom.Room.Name, orderService.Quantity, orderService.Service.Name)

		notification := model.NewNotification(id, &deptID, &bookingRoomID, model.NotificationReceiverTypeStaff, model.NotificationContentTypeOrderService, strconv.FormatInt(orderServiceID, 10), notificationContentVI, notificationContentEN, false)

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
			"code":         strconv.Itoa(constants.CodePushNotificationOrderServiceCancelled),
			"content_key":  strconv.FormatInt(orderServiceID, 10),
			"content_type": string(model.NotificationContentTypeOrderService),
		}

		u.pushNotification(schema, tokens, pnData)

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (u *orderServiceUseCaseImpl) ListForAdmin(ctx context.Context, schema string, userID int64, departmentID *int64, query request.OrderServicePaginationQuery) ([]*response.BasicOrderServiceResponse, *response.MetaResponse, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	orderServiceQuery := mapper.ToOrderServiceQuery(query)
	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return nil, nil, err
	}
	if !scope.IsGlobal() {
		orderServiceQuery.BranchIDs = scope.BranchIDs
	}

	orderServices, total, err := u.orderServiceRepo.FindAllWithServicePaginated(ctx, schema, departmentID, orderServiceQuery)
	if err != nil {
		u.logPro.Error("Find all order service paginated failed", port.String("schema", schema), port.Error(err))
		return nil, nil, err
	}

	meta := mapper.ToMetaResponse(total, query.Page, query.Limit)

	return mapper.ToBasicOrderServicesResponse(orderServices), meta, nil
}

func (u *orderServiceUseCaseImpl) UpdateStatusForAdmin(ctx context.Context, schema string, orderServiceID, userID int64, req request.UpdateOrderServiceStatusForAdminRequest) error {
	orderService, err := u.orderServiceRepo.FindByIDWithDetails(ctx, schema, orderServiceID)
	if err != nil {
		u.logPro.Error("Find order service by id failed", port.String("schema", schema), port.Int64("order_service_id", orderServiceID), port.Error(err))
		return err
	}
	if orderService == nil {
		return exceptions.ErrOrderServiceNotFound
	}
	if err := u.ensureBookingRoomInUserScope(ctx, schema, userID, orderService.BookingRoom); err != nil {
		return err
	}

	if orderService.Status != model.OrderServiceStatusPending || req.Status != model.OrderServiceStatusConfirmed {
		return exceptions.ErrInvalidStatus
	}

	if err = u.txPro.WithinTx(ctx, func(ctx context.Context) error {
		updateData := map[string]any{
			"status":        model.OrderServiceStatusConfirmed,
			"updated_by_id": userID,
		}
		if req.StaffNote != nil {
			updateData["staff_note"] = req.StaffNote
		}

		if err := u.orderServiceRepo.Update(ctx, schema, orderServiceID, updateData); err != nil {
			if errors.Is(err, exceptions.ErrOrderServiceNotFound) {
				return err
			}
			u.logPro.Error("Update order service failed", port.String("schema", schema), port.Int64("order_service_id", orderServiceID), port.Error(err))
			return err
		}

		id, err := u.idPro.NextID()
		if err != nil {
			u.logPro.Error("Generate notification id failed", port.String("schema", schema), port.Error(err))
			return err
		}

		notificationContentVI := fmt.Sprintf("%d %s đã được xác nhận", orderService.Quantity, orderService.Service.Name)
		notificationContentEN := fmt.Sprintf("%d %s have been confirmed", orderService.Quantity, orderService.Service.Name)

		notification := model.NewNotification(id, nil, &orderService.BookingRoomID, model.NotificationReceiverTypeGuest, model.NotificationContentTypeOrderService, strconv.FormatInt(orderServiceID, 10), notificationContentVI, notificationContentEN, false)

		if err = u.notificationRepo.Create(ctx, schema, notification); err != nil {
			u.logPro.Error("Create notification failed", port.String("schema", schema), port.Error(err))
			return err
		}

		deviceTokens, err := u.deviceTokenRepo.FindAllByBookingRoomID(ctx, schema, orderService.BookingRoomID)
		if err != nil {
			u.logPro.Error("Find all device token by booking room id failed", port.String("schema", schema), port.Error(err))
			return err
		}

		tokens := make([]string, 0, len(deviceTokens))
		for _, deviceToken := range deviceTokens {
			tokens = append(tokens, deviceToken.Token)
		}

		pnData := map[string]string{
			"code":         strconv.Itoa(constants.CodePushNotificationOrderServiceConfirmed),
			"content_key":  strconv.FormatInt(orderServiceID, 10),
			"content_type": string(model.NotificationContentTypeOrderService),
		}

		u.pushNotification(schema, tokens, pnData)

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (u *orderServiceUseCaseImpl) pushNotification(schema string, tokens []string, pnData map[string]string) {
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

func (u *orderServiceUseCaseImpl) ensureServiceInBookingRoomBranch(ctx context.Context, schema string, serviceID int64, bookingRoom *model.BookingRoom) error {
	if bookingRoom == nil || bookingRoom.Room == nil || bookingRoom.Room.Floor == nil || bookingRoom.Room.Floor.BranchID == 0 {
		return exceptions.ErrRoomNotInBookingBranch
	}

	exists, err := u.serviceRepo.ExistsByIDAndBranchID(ctx, schema, serviceID, bookingRoom.Room.Floor.BranchID)
	if err != nil {
		u.logPro.Error(
			"Check service branch failed",
			port.String("schema", schema),
			port.Int64("service_id", serviceID),
			port.Int64("branch_id", bookingRoom.Room.Floor.BranchID),
			port.Error(err),
		)
		return err
	}
	if !exists {
		return exceptions.ErrServiceNotInBranch
	}

	return nil
}

func (u *orderServiceUseCaseImpl) ensureBookingRoomInUserScope(ctx context.Context, schema string, userID int64, bookingRoom *model.BookingRoom) error {
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
