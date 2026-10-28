package servicetype

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
)

type ServiceTypeUseCase interface {
	Create(ctx context.Context, schema string, userID int64, req request.CreateServiceTypeRequest) error

	ListFullForAdmin(ctx context.Context, schema string, query request.ServiceTypePaginationQuery) ([]*response.ServiceTypeResponse, *response.MetaResponse, error)

	ListBasicForAdmin(ctx context.Context, schema string) ([]*response.BasicServiceTypeResponse, error)

	Update(ctx context.Context, schema string, svcTypeID, userID int64, req request.UpdateServiceTypeRequest) error

	Delete(ctx context.Context, schema string, svcTypeID int64) error

	ListWithActiveServices(ctx context.Context, schema string, query request.ServiceTypeWithActiveServicesQuery) ([]*response.SimpleServiceTypeResponse, error)

	ListForGuest(ctx context.Context, schema string) ([]*response.BasicServiceTypeResponse, error)

	ListServicesBySlug(ctx context.Context, schema, svcTypeSlug string, includeServiceType bool, branchID int64) (*response.BasicServiceTypeResponse, []*response.ServiceResponse, error)
}
