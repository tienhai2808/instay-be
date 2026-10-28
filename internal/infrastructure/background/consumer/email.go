package consumer

import (
	"encoding/json"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/data"
	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/pkg/utils"
)

func (c *Consumer) startEmailConsumer() {
	go c.startSendAuthEmail()
}

func (c *Consumer) startSendAuthEmail() {
	routingKey := utils.GetRoutingKey("*", port.RoutingKeyAuthEmail)
	queueName := utils.GetQueueName("all", port.QueueNameAuthEmail)

	if err := c.mqPro.ConsumeMessage(queueName, port.ExchangeEmail, routingKey, port.ExchangeTypeTopic, func(body []byte) error {
		var emailMsg data.AuthEmailMessage
		if err := json.Unmarshal(body, &emailMsg); err != nil {
			c.logPro.Error("Json unmarshal auth email message failed", port.String("schema", emailMsg.Schema), port.Error(err))
			return err
		}

		if err := c.mailerPro.AuthEmail(emailMsg.To, emailMsg.Subject, emailMsg.Otp); err != nil {
			c.logPro.Error("Send auth email failed", port.String("schema", emailMsg.Schema), port.Error(err))
			return err
		}

		return nil
	}); err != nil {
		c.logPro.Error("Start consumer send auth email failed", port.Error(err))
	}
}
