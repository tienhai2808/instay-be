package port

const (
	ExchangeTypeDirect = "direct"
	ExchangeTypeTopic  = "topic"
	ExchangeTypeFanout = "fanout"

	ExchangeEmail       = "send.email"
	QueueNameAuthEmail  = "email.auth"
	RoutingKeyAuthEmail = "email.auth"

	ExchangePushNotification   = "push.notification"
	QueueNamePushNotification  = "notification.push"
	RoutingKeyPushNotification = "notification.push"

	ExchangeFile         = "handle.file"
	QueueNameDeleteFile  = "file.delete"
	RoutingKeyDeleteFile = "file.delete"
)

type MessageQueueProvider interface {
	PublishMessage(exchange, routingKey string, body []byte) error

	ConsumeMessage(queueName, exchange, routingKey, exchangeType string, handler func([]byte) error) error
}
