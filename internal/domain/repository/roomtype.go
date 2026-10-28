package repository

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type RoomTypeRepository interface {
	Create(ctx context.Context, schema string, roomType *model.RoomType) error

	FindByID(ctx context.Context, schema string, id int64) (*model.RoomType, error)

	ExistsByIDAndBranchID(ctx context.Context, schema string, id, branchID int64) (bool, error)

	FindAllWithDetailsPaginated(ctx context.Context, schema string, query model.RoomTypeQuery) ([]*model.RoomType, int64, error)

	SearchByName(ctx context.Context, schema string, keyword string) (*model.RoomType, error)

	FindAll(ctx context.Context, schema string, query model.RoomTypeQuery) ([]*model.RoomType, error)

	Update(ctx context.Context, schema string, id int64, updateData map[string]any, branchIDs []int64) error

	Delete(ctx context.Context, schema string, id int64) error

	CountRoomByIDs(ctx context.Context, schema string, ids []int64) (map[int64]int64, error)
}
