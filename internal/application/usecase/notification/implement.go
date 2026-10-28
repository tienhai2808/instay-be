package notification

import (
	"context"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/mapper"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/internal/application/usecase/branchscope"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
)

type notificationUseCaseImpl struct {
	logPro               port.LogProvider
	notificationRepo     repository.NotificationRepository
	userNotificationRepo repository.UserNotificationRepository
	userRepo             repository.UserRepository
}

func NewNotificationUseCase(
	logPro port.LogProvider,
	notificationRepo repository.NotificationRepository,
	userNotificationRepo repository.UserNotificationRepository,
	userRepo repository.UserRepository,
) NotificationUseCase {
	return &notificationUseCaseImpl{
		logPro:               logPro,
		notificationRepo:     notificationRepo,
		userNotificationRepo: userNotificationRepo,
		userRepo:             userRepo,
	}
}

func (u *notificationUseCaseImpl) CountUnreadForAdmin(ctx context.Context, schema string, userID, deptID int64) (int64, error) {
	query, err := u.scopedNotificationQuery(ctx, schema, userID, request.NotificationPaginationQuery{})
	if err != nil {
		return 0, err
	}

	count, err := u.notificationRepo.CountUnreadByUserIDAndDepartmentID(ctx, schema, userID, deptID, query)
	if err != nil {
		u.logPro.Error("Count unread notification for admin failed", port.String("schema", schema), port.Error(err))
		return 0, err
	}

	return count, nil
}

func (u *notificationUseCaseImpl) CountUnreadForGuest(ctx context.Context, schema string, bookingRoomID int64) (int64, error) {
	count, err := u.notificationRepo.CountUnreadByBookingRoomID(ctx, schema, bookingRoomID)
	if err != nil {
		u.logPro.Error("Count unread notification for guest failed", port.String("schema", schema), port.Error(err))
		return 0, err
	}

	return count, nil
}

func (u *notificationUseCaseImpl) ListForAdmin(ctx context.Context, schema string, userID, deptID int64, query request.NotificationPaginationQuery) ([]*response.NotificationResponse, *response.MetaResponse, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	notificationQuery, err := u.scopedNotificationQuery(ctx, schema, userID, query)
	if err != nil {
		return nil, nil, err
	}

	notifications, total, err := u.notificationRepo.FindAllByUserIDAndDepartmentIDPaginated(ctx, schema, userID, deptID, notificationQuery)
	if err != nil {
		u.logPro.Error("Find all notifications for admin paginated failed", port.String("schema", schema), port.Error(err))
		return nil, nil, err
	}

	unreadNotificationIDs := make([]int64, 0, len(notifications))
	for _, notification := range notifications {
		if !notification.IsRead {
			unreadNotificationIDs = append(unreadNotificationIDs, notification.ID)
		}
	}

	if len(unreadNotificationIDs) > 0 {
		u.updateSeenNotifications(schema, userID, unreadNotificationIDs)
	}

	meta := mapper.ToMetaResponse(total, query.Page, query.Limit)

	return mapper.ToNotificationsResponse(notifications), meta, nil
}

func (u *notificationUseCaseImpl) ListForGuest(ctx context.Context, schema string, bookingRoomID int64, query request.NotificationPaginationQuery) ([]*response.NotificationResponse, *response.MetaResponse, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	notifications, total, err := u.notificationRepo.FindAllByBookingRoomIDPaginated(ctx, schema, bookingRoomID, mapper.ToNotificationQuery(query))
	if err != nil {
		u.logPro.Error("Find all notifications for guest paginated failed", port.String("schema", schema), port.Error(err))
		return nil, nil, err
	}

	unreadNotificationIDs := make([]int64, 0, len(notifications))
	for _, notification := range notifications {
		if !notification.IsRead {
			unreadNotificationIDs = append(unreadNotificationIDs, notification.ID)
		}
	}

	if len(unreadNotificationIDs) > 0 {
		u.updateSeenGuestNotifications(schema, bookingRoomID, unreadNotificationIDs)
	}

	meta := mapper.ToMetaResponse(total, query.Page, query.Limit)

	return mapper.ToNotificationsResponse(notifications), meta, nil
}

func (u *notificationUseCaseImpl) updateSeenNotifications(schema string, userID int64, unreadNotificationIDs []int64) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		userNotifications := make([]*model.UserNotification, 0, len(unreadNotificationIDs))
		for _, notificationID := range unreadNotificationIDs {
			userNotifications = append(userNotifications, &model.UserNotification{
				NotificationID: notificationID,
				UserID:         userID,
			})
		}

		if err := u.userNotificationRepo.CreateBatch(ctx, schema, userNotifications); err != nil {
			u.logPro.Error("Create user notifications failed", port.String("schema", schema), port.Error(err))
		}
	}()
}

func (u *notificationUseCaseImpl) updateSeenGuestNotifications(schema string, bookingRoomID int64, unreadNotificationIDs []int64) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := u.notificationRepo.UpdateReadByBookingRoomIDAndIDs(ctx, schema, bookingRoomID, unreadNotificationIDs); err != nil {
			u.logPro.Error("Update guest notifications read status failed", port.String("schema", schema), port.Error(err))
		}
	}()
}

func (u *notificationUseCaseImpl) scopedNotificationQuery(ctx context.Context, schema string, userID int64, query request.NotificationPaginationQuery) (model.NotificationQuery, error) {
	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return model.NotificationQuery{}, err
	}

	notificationQuery := mapper.ToNotificationQuery(query)
	if !scope.IsGlobal() {
		notificationQuery.BranchIDs = scope.BranchIDs
	}

	return notificationQuery, nil
}
