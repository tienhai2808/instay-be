package review

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
)

type ReviewUseCase interface {
	Create(ctx context.Context, schema string, bookingRoomID int64, req request.CreateReviewRequest) error

	ListForGuest(ctx context.Context, schema string, query request.ReviewPaginationQuery) ([]*response.SimpleReviewResponse, *response.MetaResponse, error)

	GetByBookingRoomID(ctx context.Context, schema string, bookingRoomID int64) (*response.BasicReviewResponse, error)

	ListForAdmin(ctx context.Context, schema string, userID int64, query request.ReviewPaginationQuery) ([]*response.ReviewResponse, *response.MetaResponse, error)

	UpdateByBookingRoomID(ctx context.Context, schema string, bookingRoomID int64, req request.UpdateReviewRequest) error

	Summary(ctx context.Context, schema string, userID int64) (*response.ReviewSummaryResponse, error)
}
