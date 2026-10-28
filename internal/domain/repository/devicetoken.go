package repository

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type DeviceTokenRepository interface {
	FindAllByUserIDs(ctx context.Context, schema string, userIDs []int64) ([]*model.DeviceToken, error)

	FindAllByBookingRoomID(ctx context.Context, schema string, bookingRoomID int64) ([]*model.DeviceToken, error)

	FindAllByBookingRoomIDs(ctx context.Context, schema string, bookingRoomIDs []int64) ([]*model.DeviceToken, error)

	ExistsByToken(ctx context.Context, schema string, token string) (bool, error)

	Update(ctx context.Context, schema, token string, updateData map[string]any) error

	Create(ctx context.Context, schema string, deviceToken *model.DeviceToken) error

	Delete(ctx context.Context, schema string, token string) error

	Updates(ctx context.Context, schema string, tokens []string, updateData map[string]any) error
}
