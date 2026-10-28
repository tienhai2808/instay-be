package initialization

import (
	"context"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
)

type PushNotification struct {
	fcm *messaging.Client
}

func InitPushNotification() (*PushNotification, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return nil, err
	}

	fcm, err := app.Messaging(ctx)
	if err != nil {
		return nil, err
	}

	return &PushNotification{fcm}, nil
}

func (p *PushNotification) Client() *messaging.Client {
	return p.fcm
}
