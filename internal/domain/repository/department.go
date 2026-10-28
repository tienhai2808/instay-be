package repository

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type DepartmentRepository interface {
	Create(ctx context.Context, schema string, dept *model.Department) error

	FindByIDWithPermissions(ctx context.Context, schema string, id int64) (*model.Department, error)

	UpdatePermissions(ctx context.Context, schema string, dept *model.Department, perms []*model.Permission) error

	Update(ctx context.Context, schema string, id int64, updateData map[string]any) error

	Delete(ctx context.Context, schema string, id int64) error

	FindAllByPermissionNameWithActiveUsers(ctx context.Context, schema string, permName model.PermissionName) ([]*model.Department, error)

	FindAll(ctx context.Context, schema string) ([]*model.Department, error)

	CountUserByIDs(ctx context.Context, schema string, ids []int64) (map[int64]int64, error)

	FindByName(ctx context.Context, schema, name string) (*model.Department, error)

	FindByID(ctx context.Context, schema string, id int64) (*model.Department, error)

	FindAllWithDetailsPaginated(ctx context.Context, schema string, query model.DepartmentQuery) ([]*model.Department, int64, error)
}
