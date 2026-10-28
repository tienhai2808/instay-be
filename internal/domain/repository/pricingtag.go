package repository

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type PricingTagRepository interface {
	Create(ctx context.Context, schema string, pricingTag *model.PricingTag) error

	FindAll(ctx context.Context, schema string) ([]*model.PricingTag, error)

	FindAllWithDetailsPaginated(ctx context.Context, schema string, query model.PricingTagQuery) ([]*model.PricingTag, int64, error)

	Update(ctx context.Context, schema string, id int64, updateData map[string]any) error

	Delete(ctx context.Context, schema string, id int64) error
}
