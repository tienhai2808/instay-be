package consumer

import "github.com/InstaySystem/is_v2-be/internal/application/port"

type Consumer struct {
	logPro     port.LogProvider
	mqPro      port.MessageQueueProvider
	mailerPro  port.MailerProvider
	pnPro      port.PushNotificationProvider
	storagePro port.StorageProvider
}

func NewConsumer(
	logPro port.LogProvider,
	mqPro port.MessageQueueProvider,
	mailerPro port.MailerProvider,
	pnPro port.PushNotificationProvider,
	storagePro port.StorageProvider,
) *Consumer {
	return &Consumer{
		logPro,
		mqPro,
		mailerPro,
		pnPro,
		storagePro,
	}
}

func (c *Consumer) Start() {
	c.startNotificationConsumer()
	c.startEmailConsumer()
	c.startFileConsumer()
}
