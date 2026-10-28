package repository

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type NotificationRepository interface {
	Create(ctx context.Context, schema string, notification *model.Notification) error

	CountUnreadByUserIDAndDepartmentID(ctx context.Context, schema string, userID, departmentID int64, query model.NotificationQuery) (int64, error)

	CountUnreadByBookingRoomID(ctx context.Context, schema string, bookingRoomID int64) (int64, error)

	FindAllByUserIDAndDepartmentIDPaginated(ctx context.Context, schema string, userID, departmentID int64, query model.NotificationQuery) ([]*model.Notification, int64, error)

	FindAllByBookingRoomIDPaginated(ctx context.Context, schema string, bookingRoomID int64, query model.NotificationQuery) ([]*model.Notification, int64, error)

	UpdateReadByBookingRoomIDAndIDs(ctx context.Context, schema string, bookingRoomID int64, ids []int64) error
}
