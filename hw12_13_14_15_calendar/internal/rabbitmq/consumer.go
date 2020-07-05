package rabbitmq

import (
	"net"

	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/app/sender"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/streadway/amqp"
)

type Consumer struct {
	rmq *Rabbit
}

func NewConsumer() (sender.ConsumerInterface, error) {
	rb, err := NewRabbit("amqp://"+net.JoinHostPort(config.SchedConf.Rabbit.Host, config.SchedConf.Rabbit.Port), "")
	return &Consumer{
		rmq: rb,
	}, err
}

func (c *Consumer) Receive() (<-chan amqp.Delivery, error) {
	return c.rmq.GetChan().Consume(
		c.rmq.GetQueueName(),
		"",
		false,
		false,
		false,
		false,
		nil,
	)
}

func (c *Consumer) Connect() error {
	return c.rmq.Connect()
}

func (c *Consumer) Close() error {
	return c.rmq.Close()
}
