package supportticket

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
)

type SupportTicketUseCase interface {
	Create(ctx context.Context, schema string, bookingRoomID int64, req request.CreateSupportTicketRequest) (int64, error)

	ListForGuest(ctx context.Context, schema string, bookingRoomID int64) ([]*response.SimpleSupportTicketResponse, error)

	UpdateStatusForGuest(ctx context.Context, schema string, supportTicketID, bookingRoomID int64, req request.UpdateSupportTicketStatusForGuestRequest) error

	UpdateStatusForAdmin(ctx context.Context, schema string, supportTicketID, userID int64, req request.UpdateSupportTicketStatusForAdminRequest) error

	ListForAdmin(ctx context.Context, schema string, userID int64, departmentID *int64, query request.SupportTicketPaginationQuery) ([]*response.SupportTicketResponse, *response.MetaResponse, error)
}
