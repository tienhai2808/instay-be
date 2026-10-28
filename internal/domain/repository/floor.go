package repository

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type FloorRepository interface {
	Create(ctx context.Context, schema string, room *model.Floor) error

	FindByID(ctx context.Context, schema string, id int64) (*model.Floor, error)

	FindAll(ctx context.Context, schema string, query model.FloorQuery) ([]*model.Floor, error)

	FindAllWithDetailsPaginated(ctx context.Context, schema string, query model.FloorQuery) ([]*model.Floor, int64, error)

	Update(ctx context.Context, schema string, id int64, updateData map[string]any) error

	Delete(ctx context.Context, schema string, id int64) error
}
