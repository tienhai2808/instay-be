package repository

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type ReviewRepository interface {
	Create(ctx context.Context, schema string, review *model.Review) error

	Update(ctx context.Context, schema string, id int64, updateData map[string]any) error

	FindAllPaginated(ctx context.Context, schema string, query model.ReviewQuery) ([]*model.Review, int64, error)

	FindAllWithBookingRoomDetailsPaginated(ctx context.Context, schema string, query model.ReviewQuery) ([]*model.Review, int64, error)

	CountStarSummary(ctx context.Context, schema string, branchIDs []int64) (*model.ReviewStarSummary, error)

	FindDashboardReviewSummary(ctx context.Context, schema string, branchIDs []int64) (*model.DashboardReviewSummary, error)

	FindDashboardReviewCategoryCount(ctx context.Context, schema string, branchIDs []int64) (*model.DashboardReviewCategoryCount, error)

	FindDashboardReviewCategorySentiment(ctx context.Context, schema string, branchIDs []int64) (*model.DashboardReviewCategorySentiment, error)

	FindDashboardReviewSentimentCount(ctx context.Context, schema string, branchIDs []int64) (*model.DashboardReviewSentimentCount, error)

	FindDashboardReviewSentimentTime(ctx context.Context, schema string, branchIDs []int64) (*model.DashboardReviewSentimentTime, error)

	FindByBookingRoomID(ctx context.Context, schema string, bookingRoomID int64) (*model.Review, error)

	UpdateByBookingRoomID(ctx context.Context, schema string, bookingRoomID int64, updateData map[string]any) error
}
