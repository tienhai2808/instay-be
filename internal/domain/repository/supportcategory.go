package repository

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type SupportCategoryRepository interface {
	Create(ctx context.Context, schema string, supportCategory *model.SupportCategory) error

	Update(ctx context.Context, schema string, id int64, updateData map[string]any) error

	FindAllWithDetailsPaginated(ctx context.Context, schema string, query model.SupportCategoryQuery) ([]*model.SupportCategory, int64, error)

	Delete(ctx context.Context, schema string, id int64) error

	FindAll(ctx context.Context, schema string) ([]*model.SupportCategory, error)

	FindByID(ctx context.Context, schema string, id int64) (*model.SupportCategory, error)

	CountSupportTicketByIDs(ctx context.Context, schema string, ids []int64) (map[int64]int64, error)
}
