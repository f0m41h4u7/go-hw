package sender

import (
	"context"

	"github.com/streadway/amqp"
)

type ConsumerInterface interface {
	Connect() error
	Receive(func(<-chan amqp.Delivery), context.Context) error
	Reconnect(context.Context) (<-chan amqp.Delivery, error)
	Close() error
}
