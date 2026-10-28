package repository

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type BookingRoomRepository interface {
	Create(ctx context.Context, schema string, bookingRoom *model.BookingRoom) error

	CreateBatch(ctx context.Context, schema string, bookingRooms []*model.BookingRoom) error

	Update(ctx context.Context, schema string, id int64, updateData map[string]any) error

	Updates(ctx context.Context, schema string, ids []int64, updateData map[string]any) error

	DeleteAllByBookingID(ctx context.Context, schema string, bookingID int64) error

	FindByIDWithRoom(ctx context.Context, schema string, id int64) (*model.BookingRoom, error)

	FindByIDWithBooking(ctx context.Context, schema string, id int64) (*model.BookingRoom, error)

	FindByIDWithBookingAndRoom(ctx context.Context, schema string, id int64) (*model.BookingRoom, error)

	ExistsActiveByRoomIDExceptID(ctx context.Context, schema string, id, roomID int64) (bool, error)

	ExistsActiveByRoomID(ctx context.Context, schema string, roomID int64) (bool, error)

	FindActiveRoomIDsByRoomIDs(ctx context.Context, schema string, roomIDs []int64) ([]int64, error)

	FindByIDWithDetails(ctx context.Context, schema string, id int64) (*model.BookingRoom, error)

	FindAllActiveIDs(ctx context.Context, schema string) ([]int64, error)

	FindAllActiveIDsByBranchIDs(ctx context.Context, schema string, branchIDs []int64) ([]int64, error)
}
