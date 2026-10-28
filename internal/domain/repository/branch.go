package repository

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type BranchRepository interface {
	Create(ctx context.Context, schema string, branch *model.Branch) error

	FindAll(ctx context.Context, schema string) ([]*model.Branch, error)

	FindAllWithDetailsPaginated(ctx context.Context, schema string, query model.BranchQuery) ([]*model.Branch, int64, error)

	CountByIDs(ctx context.Context, schema string, ids []int64) (int64, error)

	Update(ctx context.Context, schema string, id int64, updateData map[string]any) error

	Delete(ctx context.Context, schema string, id int64) error
}
