package repository

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type ServiceTypeRepository interface {
	Create(ctx context.Context, schema string, serviceType *model.ServiceType) error

	FindAllWithDetailsPaginated(ctx context.Context, schema string, query model.ServiceTypeQuery) ([]*model.ServiceType, int64, error)

	CountServiceByIDs(ctx context.Context, schema string, ids []int64) (map[int64]int64, error)

	Update(ctx context.Context, schema string, id int64, updateData map[string]any) error

	Delete(ctx context.Context, schema string, id int64) error

	FindAll(ctx context.Context, schema string) ([]*model.ServiceType, error)

	FindAllWithActiveServicesDetailsFiltered(ctx context.Context, schema string, query model.ServiceTypeWithServicesQuery) ([]*model.ServiceType, error)

	FindBySlug(ctx context.Context, schema, slug string) (*model.ServiceType, error)
}
