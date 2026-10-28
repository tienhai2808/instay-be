package repository

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type RoomRepository interface {
	Create(ctx context.Context, schema string, room *model.Room) error

	FindByID(ctx context.Context, schema string, id int64) (*model.Room, error)

	FindByIDWithFloor(ctx context.Context, schema string, id int64) (*model.Room, error)

	FindAllByRoomTypeID(ctx context.Context, schema string, roomTypeID int64) ([]*model.Room, error)

	ExistsByIDs(ctx context.Context, schema string, ids []int64) ([]int64, error)

	FindAllWithDetailsPaginated(ctx context.Context, schema string, query model.RoomQuery) ([]*model.Room, int64, error)

	CountStatusSummary(ctx context.Context, schema string, query model.RoomStatusQuery) (*model.RoomStatusSummary, error)

	FindTimeline(ctx context.Context, schema string, query model.RoomTimelineQuery) (*model.RoomTimeline, error)

	Update(ctx context.Context, schema string, id int64, updateData map[string]any) error

	Delete(ctx context.Context, schema string, id int64) error

	DeleteByIDs(ctx context.Context, schema string, ids []int64) (int64, error)
}
