package event

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"regexp"
	"strconv"
	"strings"
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

var eventContentImageKeyPattern = regexp.MustCompile(`data-key\s*=\s*["']([^"']+)["']`)

type eventUseCaseImpl struct {
	logPro           port.LogProvider
	idPro            port.IDProvider
	txPro            port.TxProvider
	mqPro            port.MessageQueueProvider
	eventRepo        repository.EventRepository
	branchRepo       repository.BranchRepository
	userRepo         repository.UserRepository
	notificationRepo repository.NotificationRepository
	bookingRoomRepo  repository.BookingRoomRepository
	deviceTokenRepo  repository.DeviceTokenRepository
}

func NewEventUseCase(
	logPro port.LogProvider,
	idPro port.IDProvider,
	txPro port.TxProvider,
	mqPro port.MessageQueueProvider,
	eventRepo repository.EventRepository,
	branchRepo repository.BranchRepository,
	userRepo repository.UserRepository,
	notificationRepo repository.NotificationRepository,
	bookingRoomRepo repository.BookingRoomRepository,
	deviceTokenRepo repository.DeviceTokenRepository,
) EventUseCase {
	return &eventUseCaseImpl{
		logPro,
		idPro,
		txPro,
		mqPro,
		eventRepo,
		branchRepo,
		userRepo,
		notificationRepo,
		bookingRoomRepo,
		deviceTokenRepo,
	}
}

func (u *eventUseCaseImpl) Create(ctx context.Context, schema string, userID int64, req request.CreateEventRequest) (int64, error) {
	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return 0, err
	}

	branchIDs := uniqueInt64s(req.BranchIDs)
	if err = u.ensureBranchesExist(ctx, schema, branchIDs); err != nil {
		return 0, err
	}
	if err = ensureCanManageEventBranches(scope, branchIDs); err != nil {
		return 0, err
	}

	id, err := u.idPro.NextID()
	if err != nil {
		u.logPro.Error("Generate event id failed", port.String("schema", schema), port.Error(err))
		return 0, err
	}

	eventSlug := utils.GenerateSlug(req.Title)
	var publishedAt *time.Time
	shouldSendNotification := *req.IsActive && *req.SendNotification
	if *req.IsActive {
		now := time.Now()
		publishedAt = &now
	}

	event := model.NewEvent(
		id,
		req.Thumbnail,
		req.Title,
		eventSlug,
		req.Content,
		req.Excerpt,
		*req.IsActive,
		publishedAt,
		shouldSendNotification,
		&userID,
		&userID,
	)
	event.Branches = branchesFromIDs(branchIDs)

	pushTokens := make([]string, 0)
	if err = u.txPro.WithinTx(ctx, func(txCtx context.Context) error {
		if err = u.eventRepo.Create(txCtx, schema, event); err != nil {
			if errors.Is(err, exceptions.ErrEventAlreadyExists) {
				return err
			}
			u.logPro.Error("Create event failed", port.String("schema", schema), port.Error(err))
			return err
		}

		if !shouldSendNotification {
			return nil
		}

		pushTokens, err = u.createGuestNotificationsForEvent(txCtx, schema, id, req.Title, branchIDs)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return 0, err
	}

	if shouldSendNotification {
		u.pushNotification(schema, pushTokens, map[string]string{
			"code":         strconv.Itoa(constants.CodePushNotificationEventPublished),
			"content_key":  eventSlug,
			"content_type": string(model.NotificationContentTypeEvent),
		})
	}

	return id, nil
}

func (u *eventUseCaseImpl) ListForGuest(ctx context.Context, schema string, bookingRoomID int64, query request.EventPaginationQuery) ([]*response.BasicEventResponse, *response.MetaResponse, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	isActive := true
	query.IsActive = &isActive

	branchID, err := u.bookingRoomBranchID(ctx, schema, bookingRoomID)
	if err != nil {
		return nil, nil, err
	}

	eventQuery := mapper.ToEventQuery(query)
	eventQuery.BranchID = branchID
	eventQuery.IncludeGlobal = true

	events, total, err := u.eventRepo.FindAllPaginated(ctx, schema, eventQuery)
	if err != nil {
		u.logPro.Error("Find all active event paginated failed", port.String("schema", schema), port.Error(err))
		return nil, nil, err
	}

	meta := mapper.ToMetaResponse(total, query.Page, query.Limit)

	return mapper.ToBasicEventsResponse(events), meta, nil
}

func (u *eventUseCaseImpl) GetByID(ctx context.Context, schema string, eventID, userID int64) (*response.EventDetailsResponse, error) {
	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return nil, err
	}

	event, err := u.eventRepo.FindByIDWithDetails(ctx, schema, eventID)
	if err != nil {
		u.logPro.Error("Find event by id with details failed",
			port.String("schema", schema),
			port.Int64("event_id", eventID),
			port.Error(err))
		return nil, err
	}
	if event == nil {
		return nil, exceptions.ErrEventNotFound
	}
	if err := ensureCanViewEvent(scope, event); err != nil {
		return nil, err
	}

	return mapper.ToEventDetailsResponse(event), nil
}

func (u *eventUseCaseImpl) GetBySlug(ctx context.Context, schema string, bookingRoomID int64, slug string) (*response.EventResponse, error) {
	branchID, err := u.bookingRoomBranchID(ctx, schema, bookingRoomID)
	if err != nil {
		return nil, err
	}

	event, err := u.eventRepo.FindBySlug(ctx, schema, slug)
	if err != nil {
		u.logPro.Error("Find event by slug failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}
	if event == nil {
		return nil, exceptions.ErrEventNotFound
	}
	if !eventAppliesToBranch(event, branchID) {
		return nil, exceptions.ErrEventNotFound
	}

	return mapper.ToEventResponse(event), nil
}

func (u *eventUseCaseImpl) ListForAdmin(ctx context.Context, schema string, userID int64, query request.EventPaginationQuery) ([]*response.SimpleEventResponse, *response.MetaResponse, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	eventQuery, err := u.scopedEventQuery(ctx, schema, userID, query)
	if err != nil {
		return nil, nil, err
	}

	events, total, err := u.eventRepo.FindAllPaginated(ctx, schema, eventQuery)
	if err != nil {
		u.logPro.Error("Find all event paginated failed", port.String("schema", schema), port.Error(err))
		return nil, nil, err
	}

	meta := mapper.ToMetaResponse(total, query.Page, query.Limit)

	return mapper.ToSimpleEventsResponse(events), meta, nil
}

func (u *eventUseCaseImpl) Update(ctx context.Context, schema string, eventID, userID int64, req request.UpdateEventRequest) error {
	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return err
	}

	event, err := u.eventRepo.FindByIDWithDetails(ctx, schema, eventID)
	if err != nil {
		u.logPro.Error("Find event by id with details failed",
			port.String("schema", schema),
			port.Int64("event_id", eventID),
			port.Error(err))
		return err
	}
	if event == nil {
		return exceptions.ErrEventNotFound
	}
	if err := ensureCanManageEventBranches(scope, branchIDsFromBranches(event.Branches)); err != nil {
		return err
	}

	branchIDs := uniqueInt64s(req.BranchIDs)
	if err = u.ensureBranchesExist(ctx, schema, branchIDs); err != nil {
		return err
	}
	if err = ensureCanManageEventBranches(scope, branchIDs); err != nil {
		return err
	}

	oldFileKeys := collectEventFileKeys(event.Content, event.Thumbnail)

	publishedAt := event.PublishedAt
	isNotificationSent := event.IsNotificationSent
	eventSlug := utils.GenerateSlug(req.Title)
	shouldSendNotification := publishedAt == nil &&
		*req.IsActive &&
		req.SendNotification != nil &&
		*req.SendNotification &&
		!isNotificationSent

	if publishedAt == nil && *req.IsActive {
		now := time.Now()
		publishedAt = &now
		if shouldSendNotification {
			isNotificationSent = true
		}
	}

	updateData := map[string]any{
		"title":                req.Title,
		"slug":                 eventSlug,
		"excerpt":              req.Excerpt,
		"content":              req.Content,
		"thumbnail":            req.Thumbnail,
		"is_active":            *req.IsActive,
		"published_at":         publishedAt,
		"is_notification_sent": isNotificationSent,
		"updated_by_id":        userID,
	}

	pushTokens := make([]string, 0)
	if err = u.txPro.WithinTx(ctx, func(txCtx context.Context) error {
		if err = u.eventRepo.Update(txCtx, schema, eventID, updateData, branchIDs); err != nil {
			if errors.Is(err, exceptions.ErrEventNotFound) {
				return err
			}
			if errors.Is(err, exceptions.ErrEventAlreadyExists) {
				return err
			}
			u.logPro.Error("Update event failed",
				port.String("schema", schema),
				port.Int64("event_id", eventID),
				port.Error(err))
			return err
		}

		if !shouldSendNotification {
			return nil
		}

		pushTokens, err = u.createGuestNotificationsForEvent(txCtx, schema, eventID, req.Title, branchIDs)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	if shouldSendNotification {
		u.pushNotification(schema, pushTokens, map[string]string{
			"code":         strconv.Itoa(constants.CodePushNotificationEventPublished),
			"content_key":  eventSlug,
			"content_type": string(model.NotificationContentTypeEvent),
		})
	}

	newFileKeys := collectEventFileKeys(req.Content, req.Thumbnail)
	deletedFileKeys := diffFileKeys(oldFileKeys, newFileKeys)
	if len(deletedFileKeys) > 0 {
		u.deleteFiles(schema, deletedFileKeys)
	}

	return nil
}

func (u *eventUseCaseImpl) Delete(ctx context.Context, schema string, eventID, userID int64) error {
	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return err
	}

	event, err := u.eventRepo.FindByIDWithDetails(ctx, schema, eventID)
	if err != nil {
		u.logPro.Error("Find event by id with details failed",
			port.String("schema", schema),
			port.Int64("event_id", eventID),
			port.Error(err))
		return err
	}
	if event == nil {
		return exceptions.ErrEventNotFound
	}
	if err := ensureCanManageEventBranches(scope, branchIDsFromBranches(event.Branches)); err != nil {
		return err
	}

	fileKeys := collectEventFileKeys(event.Content, event.Thumbnail)

	if err := u.eventRepo.Delete(ctx, schema, eventID); err != nil {
		if errors.Is(err, exceptions.ErrEventNotFound) {
			return err
		}
		u.logPro.Error("Delete event failed",
			port.String("schema", schema),
			port.Int64("event_id", eventID),
			port.Error(err))
		return err
	}

	if len(fileKeys) > 0 {
		u.deleteFiles(schema, fileKeys)
	}

	return nil
}

func (u *eventUseCaseImpl) DeleteBatch(ctx context.Context, schema string, userID int64, eventIDs []int64) (int64, error) {
	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return 0, err
	}

	events, err := u.eventRepo.FindByIDs(ctx, schema, eventIDs)
	if err != nil {
		u.logPro.Error("Find event by ids failed", port.String("schema", schema), port.Error(err))
		return 0, err
	}
	for _, event := range events {
		if err := ensureCanManageEventBranches(scope, branchIDsFromBranches(event.Branches)); err != nil {
			return 0, err
		}
	}

	rowsDeleted, err := u.eventRepo.DeleteByIDs(ctx, schema, eventIDs)
	if err != nil {
		u.logPro.Error("Delete event by ids failed", port.String("schema", schema), port.Error(err))
		return 0, err
	}

	fileKeys := make([]string, 0, len(events))
	for _, event := range events {
		fileKeys = append(fileKeys, collectEventFileKeys(event.Content, event.Thumbnail)...)
	}

	if len(fileKeys) > 0 {
		u.deleteFiles(schema, fileKeys)
	}

	return rowsDeleted, nil
}

func (u *eventUseCaseImpl) createGuestNotificationsForEvent(ctx context.Context, schema string, eventID int64, title string, branchIDs []int64) ([]string, error) {
	var (
		bookingRoomIDs []int64
		err            error
	)

	if len(branchIDs) == 0 {
		bookingRoomIDs, err = u.bookingRoomRepo.FindAllActiveIDs(ctx, schema)
	} else {
		bookingRoomIDs, err = u.bookingRoomRepo.FindAllActiveIDsByBranchIDs(ctx, schema, branchIDs)
	}
	if err != nil {
		u.logPro.Error("Find all active booking room ids failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}
	if len(bookingRoomIDs) == 0 {
		return []string{}, nil
	}

	contentVI := fmt.Sprintf("Sự kiện mới: %s", title)
	contentEN := fmt.Sprintf("New event: %s", title)

	for _, bookingRoomID := range bookingRoomIDs {
		notificationID, err := u.idPro.NextID()
		if err != nil {
			u.logPro.Error("Generate notification id failed", port.String("schema", schema), port.Error(err))
			return nil, err
		}

		notification := model.NewNotification(
			notificationID,
			nil,
			&bookingRoomID,
			model.NotificationReceiverTypeGuest,
			model.NotificationContentTypeEvent,
			utils.GenerateSlug(title),
			contentVI,
			contentEN,
			false,
		)

		if err = u.notificationRepo.Create(ctx, schema, notification); err != nil {
			u.logPro.Error("Create event notification failed", port.String("schema", schema), port.Error(err))
			return nil, err
		}
	}

	deviceTokens, err := u.deviceTokenRepo.FindAllByBookingRoomIDs(ctx, schema, bookingRoomIDs)
	if err != nil {
		u.logPro.Error("Find all device tokens by booking room ids failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}

	tokens := make([]string, 0, len(deviceTokens))
	for _, deviceToken := range deviceTokens {
		tokens = append(tokens, deviceToken.Token)
	}

	return uniqueNonEmptyStrings(tokens), nil
}

func (u *eventUseCaseImpl) ensureBranchesExist(ctx context.Context, schema string, branchIDs []int64) error {
	if len(branchIDs) == 0 {
		return nil
	}

	count, err := u.branchRepo.CountByIDs(ctx, schema, branchIDs)
	if err != nil {
		u.logPro.Error("Count branch by ids failed", port.String("schema", schema), port.Any("branch_ids", branchIDs), port.Error(err))
		return err
	}
	if count != int64(len(branchIDs)) {
		return exceptions.ErrBranchNotFound
	}

	return nil
}

func (u *eventUseCaseImpl) scopedEventQuery(ctx context.Context, schema string, userID int64, query request.EventPaginationQuery) (model.EventQuery, error) {
	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return model.EventQuery{}, err
	}

	branchID, branchIDs, err := scope.Apply(query.BranchID)
	if err != nil {
		return model.EventQuery{}, err
	}

	eventQuery := mapper.ToEventQuery(query)
	eventQuery.BranchID = branchID
	eventQuery.BranchIDs = branchIDs
	eventQuery.IncludeGlobal = branchID != 0 || len(branchIDs) > 0

	return eventQuery, nil
}

func (u *eventUseCaseImpl) bookingRoomBranchID(ctx context.Context, schema string, bookingRoomID int64) (int64, error) {
	bookingRoom, err := u.bookingRoomRepo.FindByIDWithRoom(ctx, schema, bookingRoomID)
	if err != nil {
		u.logPro.Error("Find booking room by id failed", port.String("schema", schema), port.Int64("booking_room_id", bookingRoomID), port.Error(err))
		return 0, err
	}
	if bookingRoom == nil || bookingRoom.Room == nil || bookingRoom.Room.Floor == nil || bookingRoom.Room.Floor.BranchID == 0 {
		return 0, exceptions.ErrBookingRoomNotFound
	}

	return bookingRoom.Room.Floor.BranchID, nil
}

func ensureCanViewEvent(scope branchscope.Scope, event *model.Event) error {
	if scope.IsGlobal() || event == nil || len(event.Branches) == 0 {
		return nil
	}

	return scope.EnsureAny(branchIDsFromBranches(event.Branches))
}

func ensureCanManageEventBranches(scope branchscope.Scope, branchIDs []int64) error {
	if scope.IsGlobal() {
		return nil
	}
	if len(branchIDs) == 0 {
		return exceptions.ErrForbidden
	}

	return scope.EnsureAll(branchIDs)
}

func eventAppliesToBranch(event *model.Event, branchID int64) bool {
	if event == nil {
		return false
	}
	if len(event.Branches) == 0 {
		return true
	}

	for _, branch := range event.Branches {
		if branch != nil && branch.ID == branchID {
			return true
		}
	}

	return false
}

func branchesFromIDs(branchIDs []int64) []*model.Branch {
	if len(branchIDs) == 0 {
		return make([]*model.Branch, 0)
	}

	branches := make([]*model.Branch, 0, len(branchIDs))
	for _, branchID := range branchIDs {
		branches = append(branches, &model.Branch{
			Base: model.Base{ID: branchID},
		})
	}

	return branches
}

func branchIDsFromBranches(branches []*model.Branch) []int64 {
	if len(branches) == 0 {
		return make([]int64, 0)
	}

	branchIDs := make([]int64, 0, len(branches))
	for _, branch := range branches {
		if branch == nil {
			continue
		}
		branchIDs = append(branchIDs, branch.ID)
	}

	return branchIDs
}

func uniqueInt64s(values []int64) []int64 {
	if len(values) == 0 {
		return make([]int64, 0)
	}

	seen := make(map[int64]struct{}, len(values))
	res := make([]int64, 0, len(values))
	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		res = append(res, value)
	}

	return res
}

func (u *eventUseCaseImpl) deleteFiles(schema string, fileKeys []string) {
	deleteFilesMsg := data.DeleteFileMessage{
		Schema: schema,
		Keys:   uniqueNonEmptyStrings(fileKeys),
	}

	if len(deleteFilesMsg.Keys) == 0 {
		return
	}

	go func(msg data.DeleteFileMessage) {
		body, err := json.Marshal(msg)
		if err != nil {
			u.logPro.Error("Json marshal failed", port.String("schema", schema), port.Error(err))
		}

		if err = u.mqPro.PublishMessage(port.ExchangeFile, port.RoutingKeyDeleteFile, body); err != nil {
			u.logPro.Error("Publish delete files message failed", port.String("schema", schema), port.Error(err))
		}
	}(deleteFilesMsg)
}

func (u *eventUseCaseImpl) pushNotification(schema string, tokens []string, pnData map[string]string) {
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

func collectEventFileKeys(content string, thumbnail *string) []string {
	keys := extractImageKeysFromEventContent(content)
	if thumbnail != nil && strings.TrimSpace(*thumbnail) != "" {
		keys = append(keys, strings.TrimSpace(*thumbnail))
	}

	return uniqueNonEmptyStrings(keys)
}

func extractImageKeysFromEventContent(content string) []string {
	matches := eventContentImageKeyPattern.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		return []string{}
	}

	keys := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		key := strings.TrimSpace(html.UnescapeString(match[1]))
		if key != "" {
			keys = append(keys, key)
		}
	}

	return uniqueNonEmptyStrings(keys)
}

func diffFileKeys(oldKeys, newKeys []string) []string {
	if len(oldKeys) == 0 {
		return []string{}
	}

	newSet := make(map[string]struct{}, len(newKeys))
	for _, key := range newKeys {
		newSet[key] = struct{}{}
	}

	deleted := make([]string, 0, len(oldKeys))
	for _, key := range oldKeys {
		if _, ok := newSet[key]; !ok {
			deleted = append(deleted, key)
		}
	}

	return uniqueNonEmptyStrings(deleted)
}

func uniqueNonEmptyStrings(values []string) []string {
	if len(values) == 0 {
		return []string{}
	}

	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		result = append(result, trimmed)
	}

	return result
}
