package pricingtag

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
)

type PricingTagUseCase interface {
	Create(ctx context.Context, schema string, userID int64, req request.CreatePricingTagRequest) error

	ListBasic(ctx context.Context, schema string) ([]*response.BasicPricingTagResponse, error)

	ListFull(ctx context.Context, schema string, query request.PricingTagPaginationQuery) ([]*response.PricingTagResponse, *response.MetaResponse, error)

	Update(ctx context.Context, schema string, pricingTagID, userID int64, req request.UpdatePricingTagRequest) error

	Delete(ctx context.Context, schema string, pricingTagID int64) error
}
