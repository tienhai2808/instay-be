package booking

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
)

type BookingUseCase interface {
	Create(ctx context.Context, schema string, userID int64, req request.CreateBookingRequest) (int64, error)

	Update(ctx context.Context, schema string, bookingID, userID int64, req request.UpdateBookingRequest) error

	UpdateStatus(ctx context.Context, schema string, bookingID, userID int64, req request.UpdateBookingStatusRequest) error

	Delete(ctx context.Context, schema string, bookingID, userID int64) error

	GetByID(ctx context.Context, schema string, bookingID, userID int64) (*response.BookingResponse, error)

	StatusSummary(ctx context.Context, schema string, userID int64) (*response.BookingStatusSummaryResponse, error)

	Timeline(ctx context.Context, schema string, userID int64, query request.BookingTimelineQuery) (*response.BookingTimelineResponse, error)

	List(ctx context.Context, schema string, userID int64, query request.BookingPaginationQuery) ([]*response.SimpleBookingResponse, *response.MetaResponse, error)
}
