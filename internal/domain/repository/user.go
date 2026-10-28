package repository

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type UserRepository interface {
	Create(ctx context.Context, schema string, user *model.User) error

	FindByUsernameWithDepartmentDetails(ctx context.Context, schema, username string) (*model.User, error)

	FindByIDWithDepartmentDetails(ctx context.Context, schema string, id int64) (*model.User, error)

	FindByIDWithDetails(ctx context.Context, schema string, id int64) (*model.User, error)

	FindByID(ctx context.Context, schema string, id int64) (*model.User, error)

	ExistsByEmail(ctx context.Context, schema, email string) (bool, error)

	FindByEmail(ctx context.Context, schema, email string) (*model.User, error)

	Update(ctx context.Context, schema string, id int64, updateData map[string]any, branchIDs ...[]int64) error

	FindAllWithDepartmentPaginated(ctx context.Context, schema string, query model.UserQuery) ([]*model.User, int64, error)

	ExistsActiveAdminExceptID(ctx context.Context, schema string, id int64) (bool, error)

	Delete(ctx context.Context, schema string, id int64) error

	DeleteByIDs(ctx context.Context, schema string, ids []int64) (int64, error)

	ExistsActiveAdmin(ctx context.Context, schema string) (bool, error)

	FindAllActiveByDepartmentID(ctx context.Context, schema string, deptID int64) ([]*model.User, error)

	FindByIDs(ctx context.Context, schema string, ids []int64) ([]*model.User, error)
}
