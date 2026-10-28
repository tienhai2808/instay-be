package repository

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type ServiceRepository interface {
	Create(ctx context.Context, schema string, service *model.Service) error

	FindByIDWithDetails(ctx context.Context, schema string, id int64) (*model.Service, error)

	FindByIDWithServiceImages(ctx context.Context, schema string, id int64) (*model.Service, error)

	Update(ctx context.Context, schema string, id int64, updateData map[string]any, branchIDs []int64) error

	Delete(ctx context.Context, schema string, id int64) error

	FindBySlugWithServiceTypeAndServiceImages(ctx context.Context, schema string, slug string) (*model.Service, error)

	DeleteByIDs(ctx context.Context, schema string, ids []int64) (int64, error)

	FindByIDsWithServiceImages(ctx context.Context, schema string, ids []int64) ([]*model.Service, error)

	FindByIDWithServiceType(ctx context.Context, schema string, id int64) (*model.Service, error)

	ExistsByIDAndBranchID(ctx context.Context, schema string, id, branchID int64) (bool, error)

	FindAllWithServiceTypePaginated(ctx context.Context, schema string, query model.ServiceQuery) ([]*model.Service, int64, error)

	FindAllWithServiceTypeAndThumbnailPaginated(ctx context.Context, schema string, query model.ServiceQuery) ([]*model.Service, int64, error)

	FindAllActiveByServiceTypeSlugWithThumbnail(ctx context.Context, schema string, serviceTypeSlug string, branchID int64) ([]*model.Service, error)
}
