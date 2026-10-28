package implement

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/mapper"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm/entity"
)

type departmentPermissionRepositoryImpl struct {
	BaseRepository
}

func NewDepartmentPermissionRepository(baseRepo BaseRepository) repository.DepartmentPermissionRepository {
	return &departmentPermissionRepositoryImpl{baseRepo}
}

func (r *departmentPermissionRepositoryImpl) CreateBatch(ctx context.Context, schema string, dps []*model.DepartmentPermission) error {
	dpEntities := mapper.ToDepartmentPermissionEntities(dps)
	return r.getDB(ctx).
		Table(r.schemaTable(schema, entity.DepartmentPermissionTableName)).
		Create(dpEntities).Error
}
