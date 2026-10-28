package fcm

import (
	"context"

	"firebase.google.com/go/v4/messaging"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/data"
	"github.com/InstaySystem/is_v2-be/internal/application/port"
)

type pushNotificationProviderImpl struct {
	client *messaging.Client
}

func NewPushNotificationProvider(client *messaging.Client) port.PushNotificationProvider {
	return &pushNotificationProviderImpl{client}
}

func (p *pushNotificationProviderImpl) Send(ctx context.Context, token string, data data.PushNotificationData) error {
	msg := &messaging.Message{
		Token: token,
		Data:  data.Data,
	}

	_, err := p.client.Send(ctx, msg)

	return err
}

func (p *pushNotificationProviderImpl) SendBatch(ctx context.Context, tokens []string, data data.PushNotificationData) error {
	msg := &messaging.MulticastMessage{
		Tokens: tokens,
		Data:   data.Data,
	}

	_, err := p.client.SendEachForMulticast(ctx, msg)

	return err
}
