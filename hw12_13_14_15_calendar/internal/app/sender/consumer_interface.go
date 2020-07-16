package sender

import "github.com/streadway/amqp"

type ConsumerInterface interface {
	Connect() error
	Receive(func(<-chan amqp.Delivery)) error
	Reconnect() (<-chan amqp.Delivery, error)
	Close() error
}
