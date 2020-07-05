package sender

import "github.com/streadway/amqp"

type ConsumerInterface interface {
	Connect() error
	Receive() (<-chan amqp.Delivery, error)
	Close() error
}
