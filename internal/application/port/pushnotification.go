package port

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/data"
)

type PushNotificationProvider interface {
	Send(ctx context.Context, token string, data data.PushNotificationData) error

	SendBatch(ctx context.Context, tokens []string, data data.PushNotificationData) error
}
