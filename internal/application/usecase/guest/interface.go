package guest

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
)

type GuestUseCase interface {
	Create(ctx context.Context, schema string, userID int64, req request.CreateGuestRequest) error

	CreateVerification(ctx context.Context, schema string, bookingRoomID int64, callbackURL string) (string, string, error)

	Verify(ctx context.Context, req request.VerifyGuestRequest) error

	Update(ctx context.Context, schema string, guestID, userID int64, req request.UpdateGuestRequest) error

	Delete(ctx context.Context, schema string, guestID, userID int64) error

	List(ctx context.Context, schema string, userID int64, query request.GuestPaginationQuery) ([]*response.GuestResponse, *response.MetaResponse, error)

	Export(ctx context.Context, schema string, userID int64, req request.ExportListGuestsRequest) ([]*response.GuestResponse, error)
}
