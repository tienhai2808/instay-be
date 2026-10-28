package orderservice

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
)

type OrderServiceUseCase interface {
	CreateForGuest(ctx context.Context, schema string, bookingRoomID int64, req request.CreateOrderServiceForGuestRequest) (int64, error)

	ListForGuest(ctx context.Context, schema string, bookingRoomID int64) ([]*response.SimpleOrderServiceResponse, error)

	UpdateStatusForGuest(ctx context.Context, schema string, orderServiceID, bookingRoomID int64, req request.UpdateOrderServiceStatusForGuestRequest) error

	CreateForAdmin(ctx context.Context, schema string, userID int64, req request.CreateOrderServiceForAdminRequest) (int64, error)

	ListForAdmin(ctx context.Context, schema string, userID int64, departmentID *int64, query request.OrderServicePaginationQuery) ([]*response.BasicOrderServiceResponse, *response.MetaResponse, error)

	UpdateStatusForAdmin(ctx context.Context, schema string, orderServiceID, userID int64, req request.UpdateOrderServiceStatusForAdminRequest) error

	GetByIDForAdmin(ctx context.Context, schema string, orderServiceID, userID int64) (*response.OrderServiceResponse, error)
}
