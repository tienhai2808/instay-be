package implement

import (
	"context"
	"errors"

	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/mapper"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm/entity"
	"gorm.io/gorm"
)

type tenantRepositoryImpl struct {
	BaseRepository
}

func NewTenantRepository(baseRepo BaseRepository) repository.TenantRepository {
	return &tenantRepositoryImpl{baseRepo}
}

func (r *tenantRepositoryImpl) FindByDomain(ctx context.Context, domain string) (*model.Tenant, error) {
	var tenantEntity entity.TenantEntity
	if err := r.db.WithContext(ctx).
		Where("domain = ? AND is_active = true", domain).
		First(&tenantEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return mapper.ToTenant(&tenantEntity), nil
}

func (r *tenantRepositoryImpl) FindAll(ctx context.Context) ([]*model.Tenant, error) {
	var tenantEntities []*entity.TenantEntity
	if err := r.db.WithContext(ctx).
		Find(&tenantEntities).Error; err != nil {
		return nil, err
	}

	return mapper.ToTenants(tenantEntities), nil
}

func (r *tenantRepositoryImpl) FindBySchemaName(ctx context.Context, schema string) (*model.Tenant, error) {
	var tenantEntity entity.TenantEntity
	if err := r.db.WithContext(ctx).
		Where("schema_name = ? AND is_active = true", schema).
		First(&tenantEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return mapper.ToTenant(&tenantEntity), nil
}

func (r *tenantRepositoryImpl) FindActive(ctx context.Context) ([]*model.Tenant, error) {
	var tenantEntities []*entity.TenantEntity
	if err := r.db.WithContext(ctx).
		Where("is_active = true").
		Find(&tenantEntities).Error; err != nil {
		return nil, err
	}

	return mapper.ToTenants(tenantEntities), nil
}

func (r *tenantRepositoryImpl) UpdateBySchemaName(ctx context.Context, schema string, updateData map[string]any) error {
	result := r.db.WithContext(ctx).
		Model(&entity.TenantEntity{}).
		Where("schema_name = ? AND is_active = true", schema).
		Updates(updateData)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return exceptions.ErrTenantNotFound
	}

	return nil
}
