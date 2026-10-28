package repository

import (
	"context"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type ChatRepository interface {
	Create(ctx context.Context, schema string, chat *model.Chat) error

	FindAllByBookingRoomIDWithDepartmentPaginated(ctx context.Context, schema string, bookingRoomID int64, query model.ChatQuery) ([]*model.Chat, int64, error)

	FindAllWithBookingRoomDetailsPaginated(ctx context.Context, schema string, departmentID *int64, query model.ChatQuery) ([]*model.Chat, int64, error)

	Update(ctx context.Context, schema string, id int64, updateData map[string]any) error

	FindByIDWithDepartment(ctx context.Context, schema string, id int64) (*model.Chat, error)

	FindByIDWithBookingRoomDetails(ctx context.Context, schema string, id int64) (*model.Chat, error)

	FindByBookingRoomIDAndDepartmentID(ctx context.Context, schema string, bookingRoomID, departmentID int64) (*model.Chat, error)

	FindIDsByBranchIDs(ctx context.Context, schema string, branchIDs []int64) ([]int64, error)

	FindAllExternalBySinceAndCursor(ctx context.Context, schema string, since, cursorAt *time.Time, cursorID *int64, limit uint32) ([]*model.Chat, bool, error)
}
