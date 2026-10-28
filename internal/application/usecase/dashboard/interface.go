package dashboard

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
)

type DashboardUseCase interface {
	BookingOverview(ctx context.Context, schema string, userID int64, query request.DashboardOverviewQuery) (*response.DashboardBookingOverviewResponse, error)

	BookingRevenue(ctx context.Context, schema string, userID int64, query request.DashboardRevenueQuery) (*response.DashboardBookingRevenueResponse, error)

	RoomTypeRevenue(ctx context.Context, schema string, userID int64, query request.DashboardOverviewQuery) (*response.DashboardRoomTypeRevenueResponse, error)

	SourceRevenue(ctx context.Context, schema string, userID int64, query request.DashboardOverviewQuery) (*response.DashboardSourceRevenueResponse, error)

	HotelOperation(ctx context.Context, schema string, userID int64) (*response.DashboardHotelOperationResponse, error)

	GuestNational(ctx context.Context, schema string, userID int64) (*response.DashboardGuestNationalResponse, error)

	SupportTicketTime(ctx context.Context, schema string, userID int64) (*response.DashboardSupportTicketTimeResponse, error)

	ServiceTypeCount(ctx context.Context, schema string, userID int64) (*response.DashboardServiceTypeCountResponse, error)

	ReviewSummary(ctx context.Context, schema string, userID int64) (*response.DashboardReviewSummaryResponse, error)

	ReviewCategoryCount(ctx context.Context, schema string, userID int64) (*response.DashboardReviewCategoryCountResponse, error)

	ReviewCategorySentiment(ctx context.Context, schema string, userID int64) (*response.DashboardReviewCategorySentimentResponse, error)

	ReviewSentimentCount(ctx context.Context, schema string, userID int64) (*response.DashboardReviewSentimentCountResponse, error)

	ReviewSentimentTime(ctx context.Context, schema string, userID int64) (*response.DashboardReviewSentimentTimeResponse, error)
}
