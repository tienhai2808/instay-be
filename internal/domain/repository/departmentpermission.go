package repository

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type DepartmentPermissionRepository interface {
	CreateBatch(ctx context.Context, schemna string, dps []*model.DepartmentPermission) error
}
