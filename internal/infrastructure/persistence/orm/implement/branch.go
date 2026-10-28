package implement

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/mapper"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm/entity"
	"gorm.io/gorm"
)

type branchRepositoryImpl struct {
	BaseRepository
}

func NewBranchRepository(baseRepo BaseRepository) repository.BranchRepository {
	return &branchRepositoryImpl{baseRepo}
}

func (r *branchRepositoryImpl) Create(ctx context.Context, schema string, branch *model.Branch) error {
	branchEntity := mapper.ToBranchEntity(branch)
	return r.getDB(ctx).
		Table(r.schemaTable(schema, entity.BranchTableName)).
		Create(branchEntity).Error
}

func (r *branchRepositoryImpl) FindAll(ctx context.Context, schema string) ([]*model.Branch, error) {
	var branchEntities []*entity.BranchEntity
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.BranchTableName)).
		Select("id", "name").
		Find(&branchEntities).Error; err != nil {
		return nil, err
	}

	return mapper.ToBranches(branchEntities), nil
}

func (r *branchRepositoryImpl) FindAllWithDetailsPaginated(ctx context.Context, schema string, query model.BranchQuery) ([]*model.Branch, int64, error) {
	var branchEntities []*entity.BranchEntity
	var total int64

	db := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.BranchTableName)).
		Model(&entity.BranchEntity{})

	db = r.applyFilters(db, query)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*model.Branch{}, 0, nil
	}

	db = db.Session(&gorm.Session{})

	db = db.Preload("CreatedBy", r.preloadWithSchema(schema, entity.UserTableName, "id", "username", "first_name", "last_name"))
	db = db.Preload("UpdatedBy", r.preloadWithSchema(schema, entity.UserTableName, "id", "username", "first_name", "last_name"))

	offset := (query.Page - 1) * query.Limit

	if err := db.Offset(int(offset)).
		Limit(int(query.Limit)).
		Find(&branchEntities).Error; err != nil {
		return nil, 0, err
	}

	return mapper.ToBranches(branchEntities), total, nil
}

func (r *branchRepositoryImpl) CountByIDs(ctx context.Context, schema string, ids []int64) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	var count int64
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.BranchTableName)).
		Where("id IN ?", ids).
		Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *branchRepositoryImpl) Update(ctx context.Context, schema string, id int64, updateData map[string]any) error {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.BranchTableName)).
		Model(&entity.BranchEntity{}).
		Where("id = ?", id).
		Updates(updateData)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return exceptions.ErrBranchNotFound
	}

	return nil
}

func (r *branchRepositoryImpl) Delete(ctx context.Context, schema string, id int64) error {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.BranchTableName)).
		Where("id = ?", id).
		Delete(&entity.BranchEntity{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return exceptions.ErrBranchNotFound
	}

	return nil
}

func (r *branchRepositoryImpl) applyFilters(db *gorm.DB, query model.BranchQuery) *gorm.DB {
	if query.Search != "" {
		term := "%" + query.Search + "%"
		db = db.Where("name ILIKE ? OR address ILIKE ? OR ward ILIKE ? OR city ILIKE ?", term, term, term, term)
	}

	if query.City != "" {
		db = db.Where("city ILIKE ?", "%"+query.City+"%")
	}

	if query.IsActive != nil {
		db = db.Where("is_active = ?", *query.IsActive)
	}

	return db
}
