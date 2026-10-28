package repository

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type PermissionRepository interface {
	Create(ctx context.Context, schema string, permission *model.Permission) error

	FindByIDs(ctx context.Context, schema string, ids []int64) ([]*model.Permission, error)

	FindAll(ctx context.Context, schema string) ([]*model.Permission, error)
}
