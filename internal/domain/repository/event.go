package repository

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type EventRepository interface {
	Create(ctx context.Context, schema string, event *model.Event) error

	FindByIDWithDetails(ctx context.Context, schema string, id int64) (*model.Event, error)

	FindBySlug(ctx context.Context, schema string, slug string) (*model.Event, error)

	FindByIDs(ctx context.Context, schema string, ids []int64) ([]*model.Event, error)

	FindAllPaginated(ctx context.Context, schema string, query model.EventQuery) ([]*model.Event, int64, error)

	Update(ctx context.Context, schema string, id int64, updateData map[string]any, branchIDs []int64) error

	Delete(ctx context.Context, schema string, id int64) error

	DeleteByIDs(ctx context.Context, schema string, ids []int64) (int64, error)
}
