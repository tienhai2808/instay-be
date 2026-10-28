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

type roomTypePriceRepositoryImpl struct {
	BaseRepository
}

func NewRoomTypePriceRepository(baseRepo BaseRepository) repository.RoomTypePriceRepository {
	return &roomTypePriceRepositoryImpl{baseRepo}
}

func (r *roomTypePriceRepositoryImpl) FindByID(ctx context.Context, schema string, id int64) (*model.RoomTypePrice, error) {
	var roomTypePriceEntity entity.RoomTypePriceEntity

	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.RoomTypePriceTableName)).
		Where("id = ?", id).
		First(&roomTypePriceEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return mapper.ToRoomTypePrice(&roomTypePriceEntity), nil
}

func (r *roomTypePriceRepositoryImpl) FindAllByRoomTypeID(ctx context.Context, schema string, roomTypeID int64) ([]*model.RoomTypePrice, error) {
	var roomTypePriceEntities []*entity.RoomTypePriceEntity

	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.RoomTypePriceTableName)).
		Where("room_type_id = ?", roomTypeID).
		Order("pricing_tag_id ASC").
		Find(&roomTypePriceEntities).Error; err != nil {
		return nil, err
	}

	return mapper.ToRoomTypePrices(roomTypePriceEntities), nil
}

func (r *roomTypePriceRepositoryImpl) FindByRoomTypeIDAndPricingTagID(ctx context.Context, schema string, roomTypeID, pricingTagID int64) (*model.RoomTypePrice, error) {
	var roomTypePriceEntity entity.RoomTypePriceEntity

	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.RoomTypePriceTableName)).
		Where("room_type_id = ? AND pricing_tag_id = ?", roomTypeID, pricingTagID).
		First(&roomTypePriceEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return mapper.ToRoomTypePrice(&roomTypePriceEntity), nil
}

func (r *roomTypePriceRepositoryImpl) Create(ctx context.Context, schema string, roomTypePrice *model.RoomTypePrice) error {
	roomTypePriceEntity := mapper.ToRoomTypePriceEntity(roomTypePrice)
	return r.getDB(ctx).
		Table(r.schemaTable(schema, entity.RoomTypePriceTableName)).
		Create(roomTypePriceEntity).Error
}

func (r *roomTypePriceRepositoryImpl) Update(ctx context.Context, schema string, id int64, updateData map[string]any) error {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.RoomTypePriceTableName)).
		Model(&entity.RoomTypePriceEntity{}).
		Where("id = ?", id).
		Updates(updateData)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return exceptions.ErrRoomTypePriceNotFound
	}

	return nil
}

func (r *roomTypePriceRepositoryImpl) Delete(ctx context.Context, schema string, id int64) error {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.RoomTypePriceTableName)).
		Where("id = ?", id).
		Delete(&entity.RoomTypePriceEntity{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return exceptions.ErrRoomTypePriceNotFound
	}

	return nil
}
