package notification

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
)

type NotificationUseCase interface {
	CountUnreadForAdmin(ctx context.Context, schema string, userID, deptID int64) (int64, error)

	CountUnreadForGuest(ctx context.Context, schema string, bookingRoomID int64) (int64, error)

	ListForAdmin(ctx context.Context, schema string, userID, deptID int64, query request.NotificationPaginationQuery) ([]*response.NotificationResponse, *response.MetaResponse, error)

	ListForGuest(ctx context.Context, schema string, bookingRoomID int64, query request.NotificationPaginationQuery) ([]*response.NotificationResponse, *response.MetaResponse, error)
}
