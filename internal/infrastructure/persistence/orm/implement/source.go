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

type sourceRepositoryImpl struct {
	BaseRepository
}

func NewSourceRepository(baseRepo BaseRepository) repository.SourceRepository {
	return &sourceRepositoryImpl{baseRepo}
}

func (r *sourceRepositoryImpl) Create(ctx context.Context, schema string, source *model.Source) error {
	srcEntity := mapper.ToSourceEntity(source)
	return r.getDB(ctx).
		Table(r.schemaTable(schema, entity.SourceTableName)).
		Create(srcEntity).Error
}

func (r *sourceRepositoryImpl) Update(ctx context.Context, schema string, id int64, updateData map[string]any) error {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.SourceTableName)).
		Model(&entity.SourceEntity{}).
		Where("id = ?", id).
		Updates(updateData)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return exceptions.ErrSourceNotFound
	}

	return nil
}

func (r *sourceRepositoryImpl) Delete(ctx context.Context, schema string, id int64) error {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.SourceTableName)).
		Where("id = ?", id).
		Delete(&entity.SourceEntity{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return exceptions.ErrSourceNotFound
	}

	return nil
}

func (r *sourceRepositoryImpl) SearchByName(ctx context.Context, schema string, keyword string) (*model.Source, error) {
	var sourceEntity entity.SourceEntity
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.SourceTableName)).
		Where("name ILIKE ?", "%"+keyword+"%").
		First(&sourceEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.ErrSourceNotFound
		}
	}

	return mapper.ToSource(&sourceEntity), nil
}

func (r *sourceRepositoryImpl) FindAll(ctx context.Context, schema string) ([]*model.Source, error) {
	var sourceEntities []*entity.SourceEntity
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.SourceTableName)).
		Select("id", "name").
		Find(&sourceEntities).Error; err != nil {
		return nil, err
	}

	return mapper.ToSources(sourceEntities), nil
}

func (r *sourceRepositoryImpl) FindAllWithDetailsPaginated(ctx context.Context, schema string, query model.SourceQuery) ([]*model.Source, int64, error) {
	var sourceEntities []*entity.SourceEntity
	var total int64

	db := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.SourceTableName)).
		Model(&entity.SourceEntity{})

	db = r.applyFilters(db, query)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*model.Source{}, 0, nil
	}

	db = db.Session(&gorm.Session{})

	db = db.Preload("CreatedBy", r.preloadWithSchema(schema, entity.UserTableName, "id", "username", "first_name", "last_name"))
	db = db.Preload("UpdatedBy", r.preloadWithSchema(schema, entity.UserTableName, "id", "username", "first_name", "last_name"))

	offset := (query.Page - 1) * query.Limit

	if err := db.Offset(int(offset)).
		Limit(int(query.Limit)).
		Find(&sourceEntities).Error; err != nil {
		return nil, 0, err
	}

	return mapper.ToSources(sourceEntities), total, nil
}

func (r *sourceRepositoryImpl) applyFilters(db *gorm.DB, query model.SourceQuery) *gorm.DB {
	if query.Search != "" {
		term := "%" + query.Search + "%"
		db = db.Where("name ILIKE ?", term)
	}

	return db
}
