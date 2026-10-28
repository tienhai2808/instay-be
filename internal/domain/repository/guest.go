package repository

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type GuestRepository interface {
	Create(ctx context.Context, schema string, guest *model.Guest) error

	Update(ctx context.Context, schema string, id int64, updateData map[string]any) error

	Delete(ctx context.Context, schema string, id int64) error

	FindByIDWithBookingRoomDetails(ctx context.Context, schema string, id int64) (*model.Guest, error)

	FindAllByBookingRoomID(ctx context.Context, schema string, bookingRoomID int64) ([]*model.Guest, error)

	FindAllWithDetailsPaginated(ctx context.Context, schema string, query model.GuestQuery) ([]*model.Guest, int64, error)

	FindAllWithDetails(ctx context.Context, schema string, query model.GuestQuery) ([]*model.Guest, error)

	FindDashboardGuestNational(ctx context.Context, schema string, branchIDs []int64) (*model.DashboardGuestNational, error)
}
