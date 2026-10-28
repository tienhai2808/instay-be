package roomtype

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
)

type RoomTypeUseCase interface {
	Create(ctx context.Context, schema string, userID int64, req request.CreateRoomTypeRequest) error

	ListFull(ctx context.Context, schema string, userID int64, query request.RoomTypePaginationQuery) ([]*response.RoomTypeResponse, *response.MetaResponse, error)

	ListBasic(ctx context.Context, schema string, userID int64, query request.RoomTypePaginationQuery) ([]*response.BasicRoomTypeResponse, error)

	ListRoomsByID(ctx context.Context, schema string, userID, roomTypeID int64, includeRoomType bool) (*response.BasicRoomTypeResponse, []*response.SimpleRoomResponse, error)

	CalculatePrice(ctx context.Context, schema string, req request.CalculateRoomTypeRateRequest) (*response.RoomTypePriceCalculationResponse, error)

	Update(ctx context.Context, schema string, roomTypeID, userID int64, req request.UpdateRoomTypeRequest) error

	Delete(ctx context.Context, schema string, roomTypeID, userID int64) error
}
