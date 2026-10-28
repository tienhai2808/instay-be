package service

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
)

type ServiceUseCase interface {
	Create(ctx context.Context, schema string, userID int64, req request.CreateServiceRequest) (int64, error)

	ListForAdmin(ctx context.Context, schema string, userID int64, query request.ServicePaginationQuery) ([]*response.SimpleServiceResponse, *response.MetaResponse, error)

	ListForGuest(ctx context.Context, schema string, query request.ServicePaginationQuery) ([]*response.ServiceResponse, *response.MetaResponse, error)

	GetByID(ctx context.Context, schema string, userID, serviceID int64) (*response.ServiceDetailsResponse, error)

	Update(ctx context.Context, schema string, serviceID, userID int64, req request.UpdateServiceRequest) error

	Delete(ctx context.Context, schema string, serviceID, userID int64) error

	DeleteBatch(ctx context.Context, schema string, userID int64, serviceIDs []int64) (int64, error)

	GetBySlug(ctx context.Context, schema string, serviceSlug string) (*response.ExtendedServiceResponse, error)
}
