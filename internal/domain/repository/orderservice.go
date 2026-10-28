package repository

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type OrderServiceRepository interface {
	Create(ctx context.Context, schema string, orderService *model.OrderService) error

	FindAllByBookingRoomIDWithServiceDetails(ctx context.Context, schema string, bookingRoomID int64) ([]*model.OrderService, error)

	FindByIDWithServiceDetails(ctx context.Context, schema string, id int64) (*model.OrderService, error)

	FindByIDWithDetails(ctx context.Context, schema string, id int64) (*model.OrderService, error)

	Update(ctx context.Context, schema string, id int64, updateData map[string]any) error

	FindAllWithServicePaginated(ctx context.Context, schema string, departmentID *int64, query model.OrderServiceQuery) ([]*model.OrderService, int64, error)

	FindDashboardServiceTypeCount(ctx context.Context, schema string, branchIDs []int64) (*model.DashboardServiceTypeCount, error)
}
