package implement

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/mapper"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm/entity"
)

type permissionRepositoryImpl struct {
	BaseRepository
}

func NewPermissionRepository(baseRepo BaseRepository) repository.PermissionRepository {
	return &permissionRepositoryImpl{baseRepo}
}

func (r *permissionRepositoryImpl) Create(ctx context.Context, schema string, permission *model.Permission) error {
	permEntity := mapper.ToPermissionEntity(permission)
	return r.getDB(ctx).
		Table(r.schemaTable(schema, entity.PermissionTableName)).
		Create(permEntity).Error
}

func (r *permissionRepositoryImpl) FindByIDs(ctx context.Context, schema string, ids []int64) ([]*model.Permission, error) {
	var permissionEntities []*entity.PermissionEntity
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.PermissionTableName)).
		Where("id IN ?", ids).
		Find(&permissionEntities).Error; err != nil {
		return nil, err
	}

	return mapper.ToPermissions(permissionEntities), nil
}

func (r *permissionRepositoryImpl) FindAll(ctx context.Context, schema string) ([]*model.Permission, error) {
	var permissionEntities []*entity.PermissionEntity
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.PermissionTableName)).
		Find(&permissionEntities).Error; err != nil {
		return nil, err
	}

	return mapper.ToPermissions(permissionEntities), nil
}
