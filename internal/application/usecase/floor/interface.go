package floor

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
)

type FloorUseCase interface {
	Create(ctx context.Context, schema string, userID int64, req request.CreateFloorRequest) error

	ListBasic(ctx context.Context, schema string, userID int64, query request.FloorPaginationQuery) ([]*response.BasicFloorResponse, error)

	ListFull(ctx context.Context, schema string, userID int64, query request.FloorPaginationQuery) ([]*response.FloorResponse, *response.MetaResponse, error)

	Update(ctx context.Context, schema string, floorID, userID int64, req request.UpdateFloorRequest) error

	Delete(ctx context.Context, schema string, floorID, userID int64) error
}
