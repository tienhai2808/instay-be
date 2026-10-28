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

type refreshTokenRepositoryImpl struct {
	BaseRepository
}

func NewRefreshTokenRepository(baseRepo BaseRepository) repository.RefreshTokenRepository {
	return &refreshTokenRepositoryImpl{baseRepo}
}

func (r *refreshTokenRepositoryImpl) Create(ctx context.Context, schema string, token *model.RefreshToken) error {
	tokenEntity := mapper.ToRefreshTokenEntity(token)
	return r.getDB(ctx).
		Table(r.schemaTable(schema, entity.RefreshTokenTableName)).
		Create(tokenEntity).Error
}

func (r *refreshTokenRepositoryImpl) Delete(ctx context.Context, schema, token string) error {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.RefreshTokenTableName)).
		Where("token = ?", token).
		Delete(&entity.RefreshTokenEntity{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return exceptions.ErrInvalidUser
	}

	return nil
}

func (r *refreshTokenRepositoryImpl) Update(ctx context.Context, schema, token string, updateData map[string]any) error {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.RefreshTokenTableName)).
		Model(&entity.RefreshTokenEntity{}).
		Where("token = ?", token).
		Updates(updateData)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return exceptions.ErrRefreshTokenNotFound
	}

	return nil
}

func (r *refreshTokenRepositoryImpl) FindByToken(ctx context.Context, schema, hashedToken string) (*model.RefreshToken, error) {
	var tokenEntity entity.RefreshTokenEntity
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.RefreshTokenTableName)).
		Where("token = ?", hashedToken).
		First(&tokenEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return mapper.ToRefreshToken(&tokenEntity), nil
}

func (r *refreshTokenRepositoryImpl) DeleteAllByUserID(ctx context.Context, schema string, userID int64) error {
	return r.getDB(ctx).
		Table(r.schemaTable(schema, entity.RefreshTokenTableName)).
		Where("user_id = ?", userID).
		Delete(&entity.RefreshTokenEntity{}).Error
}

func (r *refreshTokenRepositoryImpl) DeleteAllByUserIDs(ctx context.Context, schema string, userIDs []int64) error {
	return r.getDB(ctx).
		Table(r.schemaTable(schema, entity.RefreshTokenTableName)).
		Where("user_id IN ?", userIDs).
		Delete(&entity.RefreshTokenEntity{}).Error
}

func (r *refreshTokenRepositoryImpl) DeleteAllExpired(ctx context.Context, schema string) (int64, error) {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.RefreshTokenTableName)).
		Where("expires_at < CURRENT_TIMESTAMP").
		Delete(&entity.RefreshTokenEntity{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}
