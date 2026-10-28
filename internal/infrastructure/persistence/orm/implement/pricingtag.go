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

type pricingTagRepositoryImpl struct {
	BaseRepository
}

func NewPricingTagRepository(baseRepo BaseRepository) repository.PricingTagRepository {
	return &pricingTagRepositoryImpl{baseRepo}
}

func (r *pricingTagRepositoryImpl) Create(ctx context.Context, schema string, pricingTag *model.PricingTag) error {
	pricingTagEntity := mapper.ToPricingTagEntity(pricingTag)
	return r.getDB(ctx).
		Table(r.schemaTable(schema, entity.PricingTagTableName)).
		Create(pricingTagEntity).Error
}

func (r *pricingTagRepositoryImpl) FindAll(ctx context.Context, schema string) ([]*model.PricingTag, error) {
	var pricingTagEntities []*entity.PricingTagEntity
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.PricingTagTableName)).
		Select("id", "name").
		Find(&pricingTagEntities).Error; err != nil {
		return nil, err
	}

	return mapper.ToPricingTags(pricingTagEntities), nil
}

func (r *pricingTagRepositoryImpl) FindAllWithDetailsPaginated(ctx context.Context, schema string, query model.PricingTagQuery) ([]*model.PricingTag, int64, error) {
	var pricingTagEntities []*entity.PricingTagEntity
	var total int64

	db := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.PricingTagTableName)).
		Model(&entity.PricingTagEntity{})

	db = r.applyFilters(db, query)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*model.PricingTag{}, 0, nil
	}

	db = db.Session(&gorm.Session{})
	db = db.Preload("CreatedBy", r.preloadWithSchema(schema, entity.UserTableName, "id", "username", "first_name", "last_name"))
	db = db.Preload("UpdatedBy", r.preloadWithSchema(schema, entity.UserTableName, "id", "username", "first_name", "last_name"))

	offset := (query.Page - 1) * query.Limit

	if err := db.Offset(int(offset)).
		Limit(int(query.Limit)).
		Find(&pricingTagEntities).Error; err != nil {
		return nil, 0, err
	}

	return mapper.ToPricingTags(pricingTagEntities), total, nil
}

func (r *pricingTagRepositoryImpl) Update(ctx context.Context, schema string, id int64, updateData map[string]any) error {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.PricingTagTableName)).
		Model(&entity.PricingTagEntity{}).
		Where("id = ?", id).
		Updates(updateData)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return exceptions.ErrPricingTagNotFound
	}

	return nil
}

func (r *pricingTagRepositoryImpl) Delete(ctx context.Context, schema string, id int64) error {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.PricingTagTableName)).
		Where("id = ?", id).
		Delete(&entity.PricingTagEntity{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return exceptions.ErrPricingTagNotFound
	}

	return nil
}

func (r *pricingTagRepositoryImpl) applyFilters(db *gorm.DB, query model.PricingTagQuery) *gorm.DB {
	if query.Search != "" {
		term := "%" + query.Search + "%"
		db = db.Where("name ILIKE ?", term)
	}

	return db
}
