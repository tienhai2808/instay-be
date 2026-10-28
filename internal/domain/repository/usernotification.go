package repository

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type UserNotificationRepository interface {
	CreateBatch(ctx context.Context, schema string, userNotifications []*model.UserNotification) error
}
