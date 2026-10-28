package repository

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type SourceRepository interface {
	Create(ctx context.Context, schema string, source *model.Source) error

	Update(ctx context.Context, schema string, id int64, updateData map[string]any) error

	Delete(ctx context.Context, schema string, id int64) error

	SearchByName(ctx context.Context, schema string, keyword string) (*model.Source, error)

	FindAll(ctx context.Context, schema string) ([]*model.Source, error)

	FindAllWithDetailsPaginated(ctx context.Context, schema string, query model.SourceQuery) ([]*model.Source, int64, error)
}
