//nolint:gofumpt
package rabbitmq

import "github.com/streadway/amqp"

type Rabbit struct {
	address string
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

func NewRabbit(addr string, qName string) (*Rabbit, error) {
	rb := Rabbit{
		address: addr,
		queue:   qName,
	}
	err := rb.Connect()
	if err != nil {
		return nil, err
	}
	return &rb, nil
}

func (rb *Rabbit) GetChan() *amqp.Channel {
	return rb.channel
}

func (rb *Rabbit) GetQueueName() string {
	return rb.queue
}

func (rb *Rabbit) Connect() error {
	conn, err := amqp.Dial(rb.address)
	if err != nil {
		return err
	}
	rb.conn = conn
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	rb.channel = ch
	_, err = rb.channel.QueueDeclare(
		rb.queue,
		false,
		false,
		false,
		false,
		nil,
	)
	return err
}

func (rb *Rabbit) Close() error {
	err := rb.channel.Close()
	if err != nil {
		return err
	}
	return rb.conn.Close()
}
