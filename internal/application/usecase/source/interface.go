package source

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
)

type SourceUseCase interface {
	Create(ctx context.Context, schema string, userID int64, req request.CreateSourceRequest) error

	Update(ctx context.Context, schema string, sourceID, userID int64, req request.UpdateSourceRequest) error

	Delete(ctx context.Context, schema string, sourceID int64) error

	ListBasic(ctx context.Context, schema string) ([]*response.BasicSourceResponse, error)

	ListFull(ctx context.Context, schema string, query request.SourcePaginationQuery) ([]*response.SourceResponse, *response.MetaResponse, error)
}
