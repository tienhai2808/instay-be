package rabbitmq

import (
	"context"
	"fmt"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

const (
	publishTimeout = 5 * time.Second
	confirmTimeout = 5 * time.Second
	workerCount    = 5
	prefetch       = 5
	dlxSuffix      = ".dlx"
	dlqSuffix      = ".dlq"
	dlxRoutingKey  = "dead"
)

type messageQueueProviderImpl struct {
	rmq *amqp091.Connection
	log *zap.Logger
}

func NewMessageQueueProvider(
	rmq *amqp091.Connection,
	log *zap.Logger,
) port.MessageQueueProvider {
	return &messageQueueProviderImpl{
		rmq: rmq,
		log: log,
	}
}

func (p *messageQueueProviderImpl) PublishMessage(exchange, routingKey string, body []byte) error {
	ch, err := p.rmq.Channel()
	if err != nil {
		return err
	}
	defer func() {
		_ = ch.Close()
	}()

	if err := ch.Confirm(false); err != nil {
		return fmt.Errorf("enable publisher confirms: %w", err)
	}

	confirmCh := ch.NotifyPublish(make(chan amqp091.Confirmation, 1))
	returnCh := ch.NotifyReturn(make(chan amqp091.Return, 1))

	ctx, cancel := context.WithTimeout(context.Background(), publishTimeout)
	defer cancel()

	if err := ch.PublishWithContext(ctx, exchange, routingKey, true, false, amqp091.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp091.Persistent,
		Timestamp:    time.Now().UTC(),
		Body:         body,
	}); err != nil {
		return err
	}

	select {
	case ret := <-returnCh:
		return fmt.Errorf("message returned: exchange=%s routingKey=%s reply=%s", ret.Exchange, ret.RoutingKey, ret.ReplyText)
	case confirm := <-confirmCh:
		if !confirm.Ack {
			return fmt.Errorf("publish not confirmed by broker")
		}
	case <-time.After(confirmTimeout):
		return fmt.Errorf("publish confirm timeout after %s", confirmTimeout)
	}

	return nil
}

func (p *messageQueueProviderImpl) ConsumeMessage(queueName, exchange, routingKey, exchangeType string, handler func([]byte) error) error {
	ch, err := p.rmq.Channel()
	if err != nil {
		return err
	}

	dlx := exchange + dlxSuffix
	dlq := queueName + dlqSuffix

	if _, err := ch.QueueDeclare(dlq, true, false, false, false, nil); err != nil {
		return err
	}

	if err := ch.ExchangeDeclare(dlx, port.ExchangeTypeDirect, true, false, false, false, nil); err != nil {
		return err
	}

	if err := ch.QueueBind(dlq, dlxRoutingKey, dlx, false, nil); err != nil {
		return err
	}

	args := amqp091.Table{
		"x-dead-letter-exchange":    dlx,
		"x-dead-letter-routing-key": dlxRoutingKey,
		"x-message-ttl":             int32(86_400_000),
		"x-max-length":              int64(500_000),
		"x-max-length-bytes":        int64(500 * 1024 * 1024),
		"x-queue-type":              "quorum",
		"x-delivery-limit":          int32(5),
	}

	if _, err := ch.QueueDeclare(queueName, true, false, false, false, args); err != nil {
		return err
	}
	if err := ch.ExchangeDeclare(exchange, exchangeType, true, false, false, false, nil); err != nil {
		return err
	}
	if err := ch.QueueBind(queueName, routingKey, exchange, false, nil); err != nil {
		return err
	}

	if err := ch.Qos(prefetch, 0, false); err != nil {
		return err
	}

	consumerTag := fmt.Sprintf("%s-%d", queueName, time.Now().UnixNano())
	msgs, err := ch.Consume(queueName, consumerTag, false, false, false, false, nil)
	if err != nil {
		return err
	}

	closeErrCh := ch.NotifyClose(make(chan *amqp091.Error, 1))
	go func() {
		if closeErr := <-closeErrCh; closeErr != nil {
			p.log.Error("Rabbitmq channel closed", zap.Error(closeErr))
		}
	}()

	for i := range workerCount {
		go func(workerID int) {
			for msg := range msgs {
				p.processWithRetry(msg, handler, workerID)
			}
		}(i)
	}

	return nil
}

func (p *messageQueueProviderImpl) processWithRetry(msg amqp091.Delivery, handler func([]byte) error, workerID int) {
	err := handler(msg.Body)
	if err == nil {
		if ackErr := msg.Ack(false); ackErr != nil {
			p.log.Error("Ack message failed", zap.Int("workerID", workerID), zap.Error(ackErr))
		}
		return
	}

	p.log.Error("Handler failed, nacking to broker",
		zap.Int("workerID", workerID),
		zap.Uint64("deliveryTag", msg.DeliveryTag),
		zap.Error(err),
	)

	if nackErr := msg.Nack(false, true); nackErr != nil {
		p.log.Error("Nack message failed", zap.Int("workerID", workerID), zap.Error(nackErr))
	}
}
