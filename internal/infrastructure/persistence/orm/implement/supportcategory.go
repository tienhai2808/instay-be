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

type SupportTicketsCountResult struct {
	SupportCategoryID   int64 `gorm:"column:support_category_id"`
	SupportTicketsCount int64 `gorm:"column:support_tickets_count"`
}

type supportCategoryRepositoryImpl struct {
	BaseRepository
}

func NewSupportCategoryRepository(baseRepo BaseRepository) repository.SupportCategoryRepository {
	return &supportCategoryRepositoryImpl{baseRepo}
}

func (r *supportCategoryRepositoryImpl) Create(ctx context.Context, schema string, supportCategory *model.SupportCategory) error {
	supCatEntity := mapper.ToSupportCategoryEntity(supportCategory)
	return r.getDB(ctx).
		Table(r.schemaTable(schema, entity.SupportCategoryTableName)).
		Create(supCatEntity).Error
}

func (r *supportCategoryRepositoryImpl) Update(ctx context.Context, schema string, id int64, updateData map[string]any) error {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.SupportCategoryTableName)).
		Model(&entity.SupportCategoryEntity{}).
		Where("id = ?", id).
		Updates(updateData)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return exceptions.ErrSupportCategoryNotFound
	}

	return nil
}

func (r *supportCategoryRepositoryImpl) CountSupportTicketByIDs(ctx context.Context, schema string, ids []int64) (map[int64]int64, error) {
	var counts []SupportTicketsCountResult
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.SupportTicketTableName)).
		Select("support_category_id, COUNT(*) as support_tickets_count").
		Where("support_category_id IN ?", ids).
		Group("support_category_id").
		Scan(&counts).Error; err != nil {
		return nil, err
	}

	countMap := make(map[int64]int64, len(counts))
	for _, c := range counts {
		countMap[c.SupportCategoryID] = c.SupportTicketsCount
	}

	return countMap, nil
}

func (r *supportCategoryRepositoryImpl) FindAllWithDetailsPaginated(ctx context.Context, schema string, query model.SupportCategoryQuery) ([]*model.SupportCategory, int64, error) {
	var supCatEntities []*entity.SupportCategoryEntity
	var total int64

	db := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.SupportCategoryTableName)).
		Model(&entity.SupportCategoryEntity{})

	db = r.applyFilters(db, query)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*model.SupportCategory{}, 0, nil
	}

	db = db.Session(&gorm.Session{})

	db = db.Preload("CreatedBy", r.preloadWithSchema(schema, entity.UserTableName, "id", "username", "first_name", "last_name"))
	db = db.Preload("UpdatedBy", r.preloadWithSchema(schema, entity.UserTableName, "id", "username", "first_name", "last_name"))
	db = db.Preload("Department", r.preloadWithSchema(schema, entity.DepartmentTableName, "id", "name"))

	offset := (query.Page - 1) * query.Limit

	if err := db.Offset(int(offset)).
		Limit(int(query.Limit)).
		Find(&supCatEntities).Error; err != nil {
		return nil, 0, err
	}

	return mapper.ToSupportCategories(supCatEntities), total, nil
}

func (r *supportCategoryRepositoryImpl) FindAll(ctx context.Context, schema string) ([]*model.SupportCategory, error) {
	var supCatEntities []*entity.SupportCategoryEntity
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.SupportCategoryTableName)).
		Find(&supCatEntities).Error; err != nil {
		return nil, err
	}

	return mapper.ToSupportCategories(supCatEntities), nil
}

func (r *supportCategoryRepositoryImpl) FindByID(ctx context.Context, schema string, id int64) (*model.SupportCategory, error) {
	var supCatEntity entity.SupportCategoryEntity
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.SupportCategoryTableName)).
		Where("id = ?", id).
		First(&supCatEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}

	return mapper.ToSupportCategory(&supCatEntity), nil
}

func (r *supportCategoryRepositoryImpl) Delete(ctx context.Context, schema string, id int64) error {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.SupportCategoryTableName)).
		Where("id = ?", id).
		Delete(&entity.SupportCategoryEntity{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return exceptions.ErrSupportCategoryNotFound
	}

	return nil
}

func (r *supportCategoryRepositoryImpl) applyFilters(db *gorm.DB, query model.SupportCategoryQuery) *gorm.DB {
	if query.Search != "" {
		term := "%" + query.Search + "%"
		db = db.Where("name ILIKE ?", term)
	}

	if query.DepartmentID != 0 {
		db = db.Where("department_id = ?", query.DepartmentID)
	}

	return db
}
