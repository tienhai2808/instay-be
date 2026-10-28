package event

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
)

type EventUseCase interface {
	Create(ctx context.Context, schema string, userID int64, req request.CreateEventRequest) (int64, error)

	ListForGuest(ctx context.Context, schema string, bookingRoomID int64, query request.EventPaginationQuery) ([]*response.BasicEventResponse, *response.MetaResponse, error)

	GetByID(ctx context.Context, schema string, eventID, userID int64) (*response.EventDetailsResponse, error)

	GetBySlug(ctx context.Context, schema string, bookingRoomID int64, slug string) (*response.EventResponse, error)

	ListForAdmin(ctx context.Context, schema string, userID int64, query request.EventPaginationQuery) ([]*response.SimpleEventResponse, *response.MetaResponse, error)

	Update(ctx context.Context, schema string, eventID, userID int64, req request.UpdateEventRequest) error

	Delete(ctx context.Context, schema string, eventID, userID int64) error

	DeleteBatch(ctx context.Context, schema string, userID int64, eventIDs []int64) (int64, error)
}
