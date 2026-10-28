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

type floorRepositoryImpl struct {
	BaseRepository
}

func NewFloorRepository(baseRepo BaseRepository) repository.FloorRepository {
	return &floorRepositoryImpl{baseRepo}
}

func (r *floorRepositoryImpl) Create(ctx context.Context, schema string, floor *model.Floor) error {
	floorEntity := mapper.ToFloorEntity(floor)
	return r.getDB(ctx).
		Table(r.schemaTable(schema, entity.FloorTableName)).
		Create(floorEntity).Error
}

func (r *floorRepositoryImpl) FindByID(ctx context.Context, schema string, id int64) (*model.Floor, error) {
	var floorEntity entity.FloorEntity
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.FloorTableName)).
		Where("id = ?", id).
		First(&floorEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return mapper.ToFloor(&floorEntity), nil
}

func (r *floorRepositoryImpl) FindAll(ctx context.Context, schema string, query model.FloorQuery) ([]*model.Floor, error) {
	var floorEntities []*entity.FloorEntity

	db := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.FloorTableName)).
		Select("id", "name", "branch_id").
		Preload("Branch", r.preloadWithSchema(schema, entity.BranchTableName, "id", "name"))

	db = r.applyFilters(db, query)

	if err := db.Find(&floorEntities).Error; err != nil {
		return nil, err
	}

	return mapper.ToFloors(floorEntities), nil
}

func (r *floorRepositoryImpl) FindAllWithDetailsPaginated(ctx context.Context, schema string, query model.FloorQuery) ([]*model.Floor, int64, error) {
	var floorEntities []*entity.FloorEntity
	var total int64

	db := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.FloorTableName)).
		Model(&entity.FloorEntity{})

	db = r.applyFilters(db, query)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*model.Floor{}, 0, nil
	}

	db = db.Session(&gorm.Session{})

	db = db.Preload("Branch", r.preloadWithSchema(schema, entity.BranchTableName, "id", "name"))
	db = db.Preload("CreatedBy", r.preloadWithSchema(schema, entity.UserTableName, "id", "username", "first_name", "last_name"))
	db = db.Preload("UpdatedBy", r.preloadWithSchema(schema, entity.UserTableName, "id", "username", "first_name", "last_name"))

	offset := (query.Page - 1) * query.Limit

	if err := db.Offset(int(offset)).
		Limit(int(query.Limit)).
		Find(&floorEntities).Error; err != nil {
		return nil, 0, err
	}

	return mapper.ToFloors(floorEntities), total, nil
}

func (r *floorRepositoryImpl) Update(ctx context.Context, schema string, id int64, updateData map[string]any) error {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.FloorTableName)).
		Model(&entity.FloorEntity{}).
		Where("id = ?", id).
		Updates(updateData)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return exceptions.ErrFloorNotFound
	}

	return nil
}

func (r *floorRepositoryImpl) Delete(ctx context.Context, schema string, id int64) error {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.FloorTableName)).
		Where("id = ?", id).
		Delete(&entity.FloorEntity{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return exceptions.ErrFloorNotFound
	}

	return nil
}

func (r *floorRepositoryImpl) applyFilters(db *gorm.DB, query model.FloorQuery) *gorm.DB {
	if query.Search != "" {
		term := "%" + query.Search + "%"
		db = db.Where("name ILIKE ?", term)
	}

	if query.BranchID != 0 {
		db = db.Where("branch_id = ?", query.BranchID)
	} else if len(query.BranchIDs) > 0 {
		db = db.Where("branch_id IN ?", query.BranchIDs)
	}

	return db
}
