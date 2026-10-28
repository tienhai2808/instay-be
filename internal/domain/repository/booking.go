package repository

import (
	"context"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type BookingRepository interface {
	Create(ctx context.Context, schema string, booking *model.Booking) error

	Update(ctx context.Context, schema string, id int64, updateData map[string]any) error

	Delete(ctx context.Context, schema string, id int64) error

	FindByBookingNumber(ctx context.Context, schema, bookingNumber string) (*model.Booking, error)

	FindByIDWithBookingRooms(ctx context.Context, schema string, id int64) (*model.Booking, error)

	FindByIDWithSourceAndBookingRooms(ctx context.Context, schema string, id int64) (*model.Booking, error)

	FindByIDWithDetails(ctx context.Context, schema string, id int64) (*model.Booking, error)

	FindAllWithDetailsPaginated(ctx context.Context, schema string, query model.BookingQuery) ([]*model.Booking, int64, error)

	CountStatusSummary(ctx context.Context, schema string, branchIDs ...[]int64) (*model.BookingStatusSummary, error)

	FindTimeline(ctx context.Context, schema string, query model.BookingTimelineQuery) (*model.BookingTimeline, error)

	FindDashboardOverview(ctx context.Context, schema string, query model.DashboardOverviewQuery) (*model.DashboardBookingOverview, error)

	FindDashboardRevenue(ctx context.Context, schema string, query model.DashboardRevenueQuery) (*model.DashboardBookingRevenue, error)

	FindDashboardRoomTypeRevenue(ctx context.Context, schema string, query model.DashboardOverviewQuery) (*model.DashboardRoomTypeRevenue, error)

	FindDashboardSourceRevenue(ctx context.Context, schema string, query model.DashboardOverviewQuery) (*model.DashboardSourceRevenue, error)

	CountDashboardTodayCheckInOut(ctx context.Context, schema string, fromDate, toDate time.Time, branchIDs []int64) (int64, int64, error)
}
