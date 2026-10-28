package customer

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
)

type CustomerUseCase interface {
	Create(ctx context.Context, schema string, bookingRoomID int64, req request.CreateCustomerRequest) error

	List(ctx context.Context, schema string, userID int64, query request.CustomerPaginationQuery) ([]*response.CustomerResponse, *response.MetaResponse, error)

	Export(ctx context.Context, schema string, userID int64, req request.ExportListCustomersRequest) ([]*response.CustomerResponse, error)
}
