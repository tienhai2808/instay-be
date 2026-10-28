package room

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
)

type RoomUseCase interface {
	Create(ctx context.Context, schema string, userID int64, req request.CreateRoomRequest) error

	List(ctx context.Context, schema string, userID int64, query request.RoomPaginationQuery) ([]*response.RoomResponse, *response.MetaResponse, error)

	Timeline(ctx context.Context, schema string, userID int64, query request.RoomTimelineQuery) (*response.RoomTimelineResponse, error)

	Update(ctx context.Context, schema string, roomID, userID int64, req request.UpdateRoomRequest) error

	Delete(ctx context.Context, schema string, roomID, userID int64) error

	DeleteBatch(ctx context.Context, schema string, userID int64, roomIDs []int64) (int64, error)
}
