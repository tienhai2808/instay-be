package roomtypeprice

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
)

type RoomTypePriceUseCase interface {
	Create(ctx context.Context, schema string, req request.CreateRoomTypePriceRequest) error

	Update(ctx context.Context, schema string, roomTypePriceID int64, req request.UpdateRoomTypePriceRequest) error

	Delete(ctx context.Context, schema string, roomTypePriceID int64) error
}
