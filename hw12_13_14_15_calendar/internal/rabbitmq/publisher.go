package rabbitmq

import (
	"log"
	"net"

	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/app/scheduler"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/streadway/amqp"
)

type Publisher struct {
	rmq *Rabbit
}

func NewPublisher() (scheduler.PublisherInterface, error) {
	rb, err := NewRabbit("amqp://"+net.JoinHostPort(config.SchedConf.Rabbit.Host, config.SchedConf.Rabbit.Port), "")
	return &Publisher{
		rmq: rb,
	}, err
}

func (p *Publisher) Send(data []byte) error {
	log.Printf(" [x] Sent %s", data)
	return p.rmq.GetChan().Publish(
		"",
		p.rmq.GetQueueName(),
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         data,
		})
}

func (p *Publisher) Connect() error {
	return p.rmq.Connect()
}

func (p *Publisher) Close() error {
	return p.rmq.Close()
}
