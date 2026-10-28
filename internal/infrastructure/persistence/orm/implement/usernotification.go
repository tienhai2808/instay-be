package implement

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/mapper"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm/entity"
)

type userNotificationRepositoryImpl struct {
	BaseRepository
}

func NewUserNotificationRepository(baseRepo BaseRepository) repository.UserNotificationRepository {
	return &userNotificationRepositoryImpl{baseRepo}
}

func (r *userNotificationRepositoryImpl) CreateBatch(ctx context.Context, schema string, userNotifications []*model.UserNotification) error {
	unEntities := mapper.ToUserNotificationEntities(userNotifications)
	return r.getDB(ctx).
		Table(r.schemaTable(schema, entity.UserNotificationTableName)).
		Create(unEntities).Error
}
