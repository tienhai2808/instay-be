package bookingroom

import (
	"context"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
)

type BookingRoomUseCase interface {
	Create(ctx context.Context, schema string, userID int64, req request.CreateBookingRoomRequest) (int64, error)

	Verify(ctx context.Context, schema, code string) (*response.SimpleBookingRoomResponse, string, time.Duration, error)

	Info(ctx context.Context, schema string, bookingRoomID int64) (*response.SimpleBookingRoomResponse, error)

	GetByIDFull(ctx context.Context, schema string, bookingRoomID, userID int64) (*response.BookingRoomResponse, error)

	GetByIDBasic(ctx context.Context, schema string, bookingRoomID, userID int64) (*response.SimpleBookingRoomResponse, error)

	ListGuestsByID(ctx context.Context, schema string, bookingRoomID, userID int64) ([]*response.SimpleGuestResponse, error)

	Update(ctx context.Context, schema string, bookingRoomID, userID int64, req request.UpdateBookingRoomRequest) error
}
