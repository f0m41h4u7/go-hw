package sender

import (
	"context"

	"github.com/streadway/amqp"
)

type ConsumerInterface interface {
	Connect() error
	Receive(context.Context, func(<-chan amqp.Delivery)) error
	Reconnect(context.Context) (<-chan amqp.Delivery, error)
	Close() error
}
