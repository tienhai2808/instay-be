package implement

import (
	"context"
	"errors"
	"fmt"

	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/mapper"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm/entity"
	"gorm.io/gorm"
)

type ServicesCountResult struct {
	ServiceTypeID int64 `gorm:"column:service_type_id"`
	ServicesCount int64 `gorm:"column:services_count"`
}

type serviceTypeRepositoryImpl struct {
	BaseRepository
}

func NewServiceTypeRepository(baseRepo BaseRepository) repository.ServiceTypeRepository {
	return &serviceTypeRepositoryImpl{baseRepo}
}

func (r *serviceTypeRepositoryImpl) Create(ctx context.Context, schema string, serviceType *model.ServiceType) error {
	svcTypeEntity := mapper.ToServiceTypeEntity(serviceType)
	return r.getDB(ctx).
		Table(r.schemaTable(schema, entity.ServiceTypeTableName)).
		Create(svcTypeEntity).Error
}

func (r *serviceTypeRepositoryImpl) FindAllWithActiveServicesDetailsFiltered(ctx context.Context, schema string, query model.ServiceTypeWithServicesQuery) ([]*model.ServiceType, error) {
	var serviceTypeEntities []*entity.ServiceTypeEntity
	serviceTable := r.schemaTable(schema, entity.ServiceTableName)
	branchServiceTable := r.schemaTable(schema, entity.BranchServiceTableName)
	serviceImageTable := r.schemaTable(schema, entity.ServiceImageTableName)

	db := r.getDB(ctx).Table(r.schemaTable(schema, entity.ServiceTypeTableName))

	if query.ServiceTypeID != 0 {
		db = db.Where("id = ?", query.ServiceTypeID)
	}

	db = db.Preload("Services", func(db *gorm.DB) *gorm.DB {
		db = db.
			Table(serviceTable + " AS svc").
			Where("svc.is_active = true")

		if query.BranchID != 0 {
			db = db.Joins(fmt.Sprintf("JOIN %s bs ON bs.service_id = svc.id", branchServiceTable)).
				Where("bs.branch_id = ?", query.BranchID)
		}

		if query.Search != "" {
			term := "%" + query.Search + "%"
			db = db.Where("svc.name ILIKE ?", term)
		}

		return db.Preload("ServiceImages", func(db *gorm.DB) *gorm.DB {
			return db.
				Table(serviceImageTable).
				Where("is_thumbnail = true")
		})
	})

	if err := db.Find(&serviceTypeEntities).Error; err != nil {
		return nil, err
	}

	return mapper.ToServiceTypes(serviceTypeEntities), nil
}

func (r *serviceTypeRepositoryImpl) FindBySlug(ctx context.Context, schema, slug string) (*model.ServiceType, error) {
	var serviceTypeEntity entity.ServiceTypeEntity

	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.ServiceTypeTableName)).
		Where("slug = ?", slug).
		First(&serviceTypeEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return mapper.ToServiceType(&serviceTypeEntity), nil
}

func (r *serviceTypeRepositoryImpl) FindAllWithDetailsPaginated(ctx context.Context, schema string, query model.ServiceTypeQuery) ([]*model.ServiceType, int64, error) {
	var serviceTypeEntities []*entity.ServiceTypeEntity
	var total int64

	db := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.ServiceTypeTableName)).
		Model(&entity.ServiceTypeEntity{})

	db = r.applyFilters(db, query)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*model.ServiceType{}, 0, nil
	}

	db = db.Session(&gorm.Session{})

	db = db.Preload("CreatedBy", r.preloadWithSchema(schema, entity.UserTableName, "id", "username", "first_name", "last_name"))
	db = db.Preload("UpdatedBy", r.preloadWithSchema(schema, entity.UserTableName, "id", "username", "first_name", "last_name"))
	db = db.Preload("Department", r.preloadWithSchema(schema, entity.DepartmentTableName, "id", "name"))

	offset := (query.Page - 1) * query.Limit

	if err := db.Offset(int(offset)).
		Limit(int(query.Limit)).
		Find(&serviceTypeEntities).Error; err != nil {
		return nil, 0, err
	}

	return mapper.ToServiceTypes(serviceTypeEntities), total, nil
}

func (r *serviceTypeRepositoryImpl) CountServiceByIDs(ctx context.Context, schema string, ids []int64) (map[int64]int64, error) {
	var counts []ServicesCountResult
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.ServiceTableName)).
		Select("service_type_id, COUNT(*) as services_count").
		Where("service_type_id IN ?", ids).
		Group("service_type_id").
		Scan(&counts).Error; err != nil {
		return nil, err
	}

	countMap := make(map[int64]int64, len(counts))
	for _, c := range counts {
		countMap[c.ServiceTypeID] = c.ServicesCount
	}

	return countMap, nil
}

func (r *serviceTypeRepositoryImpl) Update(ctx context.Context, schema string, id int64, updateData map[string]any) error {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.ServiceTypeTableName)).
		Model(&entity.ServiceTypeEntity{}).
		Where("id = ?", id).
		Updates(updateData)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return exceptions.ErrServiceTypeNotFound
	}

	return nil
}

func (r *serviceTypeRepositoryImpl) Delete(ctx context.Context, schema string, id int64) error {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.ServiceTypeTableName)).
		Where("id = ?", id).
		Delete(&entity.ServiceTypeEntity{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return exceptions.ErrServiceTypeNotFound
	}

	return nil
}

func (r *serviceTypeRepositoryImpl) FindAll(ctx context.Context, schema string) ([]*model.ServiceType, error) {
	var serviceTypeEntities []*entity.ServiceTypeEntity
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.ServiceTypeTableName)).
		Find(&serviceTypeEntities).Error; err != nil {
		return nil, err
	}

	return mapper.ToServiceTypes(serviceTypeEntities), nil
}

func (r *serviceTypeRepositoryImpl) applyFilters(db *gorm.DB, query model.ServiceTypeQuery) *gorm.DB {
	if query.Search != "" {
		term := "%" + query.Search + "%"
		db = db.Where("name ILIKE ?", term)
	}

	if query.DepartmentID != 0 {
		db = db.Where("department_id = ?", query.DepartmentID)
	}

	return db
}
