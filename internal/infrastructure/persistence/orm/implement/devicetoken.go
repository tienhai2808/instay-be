package implement

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/mapper"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm/entity"
)

type deviceTokenRepositoryImpl struct {
	BaseRepository
}

func NewDeviceTokenRepository(baseRepo BaseRepository) repository.DeviceTokenRepository {
	return &deviceTokenRepositoryImpl{baseRepo}
}

func (r *deviceTokenRepositoryImpl) FindAllByUserIDs(ctx context.Context, schema string, userIDs []int64) ([]*model.DeviceToken, error) {
	var tokenEntities []*entity.DeviceTokenEntity
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.DeviceTokenTableName)).
		Where("client_id IN ? AND client_type = ?", userIDs, model.DeviceTokenClientTypeUser).
		Find(&tokenEntities).Error; err != nil {
		return nil, err
	}

	return mapper.ToDeviceTokens(tokenEntities), nil
}

func (r *deviceTokenRepositoryImpl) FindAllByBookingRoomID(ctx context.Context, schema string, bookingRoomID int64) ([]*model.DeviceToken, error) {
	var tokenEntities []*entity.DeviceTokenEntity
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.DeviceTokenTableName)).
		Where("client_id = ? AND client_type = ?", bookingRoomID, model.DeviceTokenClientTypeBookingRoom).
		Find(&tokenEntities).Error; err != nil {
		return nil, err
	}

	return mapper.ToDeviceTokens(tokenEntities), nil
}

func (r *deviceTokenRepositoryImpl) ExistsByToken(ctx context.Context, schema string, token string) (bool, error) {
	var count int64
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.DeviceTokenTableName)).
		Model(&entity.DeviceTokenEntity{}).
		Where("token = ?", token).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *deviceTokenRepositoryImpl) FindAllByBookingRoomIDs(ctx context.Context, schema string, bookingRoomIDs []int64) ([]*model.DeviceToken, error) {
	if len(bookingRoomIDs) == 0 {
		return make([]*model.DeviceToken, 0), nil
	}

	var tokenEntities []*entity.DeviceTokenEntity
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.DeviceTokenTableName)).
		Where("client_id IN ? AND client_type = ?", bookingRoomIDs, model.DeviceTokenClientTypeBookingRoom).
		Find(&tokenEntities).Error; err != nil {
		return nil, err
	}

	return mapper.ToDeviceTokens(tokenEntities), nil
}

func (r *deviceTokenRepositoryImpl) Update(ctx context.Context, schema, token string, updateData map[string]any) error {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.DeviceTokenTableName)).
		Model(&entity.DeviceTokenEntity{}).
		Where("token = ?", token).
		Updates(updateData)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return exceptions.ErrDeviceTokenNotFound
	}

	return nil
}

func (r *deviceTokenRepositoryImpl) Create(ctx context.Context, schema string, deviceToken *model.DeviceToken) error {
	tokenEntity := mapper.ToDeviceTokenEntity(deviceToken)
	return r.getDB(ctx).
		Table(r.schemaTable(schema, entity.DeviceTokenTableName)).
		Create(tokenEntity).Error
}

func (r *deviceTokenRepositoryImpl) Delete(ctx context.Context, schema string, token string) error {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.DeviceTokenTableName)).
		Where("token = ?", token).
		Delete(&entity.DeviceTokenEntity{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return exceptions.ErrDeviceTokenNotFound
	}

	return nil
}

func (r *deviceTokenRepositoryImpl) Updates(ctx context.Context, schema string, tokens []string, updateData map[string]any) error {
	return r.getDB(ctx).
		Table(r.schemaTable(schema, entity.DeviceTokenTableName)).
		Model(&entity.DeviceTokenEntity{}).
		Where("token IN ?", tokens).
		Updates(updateData).Error
}
