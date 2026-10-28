package repository

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type SupportTicketRepository interface {
	Create(ctx context.Context, schema string, supportTicket *model.SupportTicket) error

	FindAllByBookingRoomIDWithSupportCategory(ctx context.Context, schema string, bookingRoomID int64) ([]*model.SupportTicket, error)

	FindByIDWithSupportCategory(ctx context.Context, schema string, id int64) (*model.SupportTicket, error)

	Update(ctx context.Context, schema string, id int64, updateData map[string]any) error

	FindAllWithSupportCategoryPaginated(ctx context.Context, schema string, departmentID *int64, query model.SupportTicketQuery) ([]*model.SupportTicket, int64, error)

	CountDashboardPendingAndInProgress(ctx context.Context, schema string, branchIDs []int64) (int64, int64, error)

	FindDashboardSupportTicketTime(ctx context.Context, schema string, branchIDs []int64) (*model.DashboardSupportTicketTime, error)
}
