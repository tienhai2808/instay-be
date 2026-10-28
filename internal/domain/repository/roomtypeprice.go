package repository

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type RoomTypePriceRepository interface {
	FindByID(ctx context.Context, schema string, id int64) (*model.RoomTypePrice, error)

	FindAllByRoomTypeID(ctx context.Context, schema string, roomTypeID int64) ([]*model.RoomTypePrice, error)

	FindByRoomTypeIDAndPricingTagID(ctx context.Context, schema string, roomTypeID, pricingTagID int64) (*model.RoomTypePrice, error)

	Create(ctx context.Context, schema string, roomTypePrice *model.RoomTypePrice) error

	Update(ctx context.Context, schema string, id int64, updateData map[string]any) error

	Delete(ctx context.Context, schema string, id int64) error
}
