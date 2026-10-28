package consumer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/data"
	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/pkg/utils"
)

func (c *Consumer) startNotificationConsumer() {
	go c.startPushNotification()

}

func (c *Consumer) startPushNotification() {
	routingKey := utils.GetRoutingKey("*", port.RoutingKeyPushNotification)
	queueName := utils.GetQueueName("all", port.QueueNamePushNotification)

	if err := c.mqPro.ConsumeMessage(queueName, port.ExchangePushNotification, routingKey, port.ExchangeTypeTopic, func(body []byte) error {
		var pushNotificationMsg data.PushNotificationMessage
		if err := json.Unmarshal(body, &pushNotificationMsg); err != nil {
			c.logPro.Error("Json unmarshal push notification message failed", port.String("schema", pushNotificationMsg.Schema), port.Error(err))
			return err
		}

		pushNotificationData := data.PushNotificationData{
			Data: pushNotificationMsg.Data,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := c.pnPro.SendBatch(ctx, pushNotificationMsg.Tokens, pushNotificationData); err != nil {
			c.logPro.Error("Batch push notifications failed", port.String("schema", pushNotificationMsg.Schema), port.Error(err))
			return err
		}

		return nil
	}); err != nil {
		c.logPro.Error("Start consumer push notification failed", port.Error(err))
	}
}
