package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/cenkalti/backoff/v3"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/app/scheduler"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/streadway/amqp"
)

type Publisher struct {
	address string
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
	done    chan error
}

func NewPublisher() scheduler.PublisherInterface {
	return &Publisher{
		address: "amqp://" + net.JoinHostPort(config.SendConf.Rabbit.Host, config.SendConf.Rabbit.Port),
		queue:   "eventQueue",
		done:    make(chan error),
	}
}

func (p *Publisher) Send(data []byte) error {
	be := backoff.NewExponentialBackOff()
	be.MaxElapsedTime = time.Minute
	be.InitialInterval = 1 * time.Second
	be.Multiplier = 2
	be.MaxInterval = 15 * time.Second

	b := backoff.WithContext(be, context.Background())
	for {
		d := b.NextBackOff()
		if d == backoff.Stop {
			return fmt.Errorf("stop reconnecting")
		}

		//nolint:gosimple
		select {
		case <-time.After(d):
			if err := p.Connect(); err != nil {
				log.Printf("could not connect in reconnect call: %+v", err)
				continue
			}
			err := p.channel.Publish(
				"",
				p.queue,
				false,
				false,
				amqp.Publishing{
					DeliveryMode: amqp.Persistent,
					ContentType:  "application/json",
					Body:         data,
				})
			if err != nil {
				fmt.Printf("failed to send data: %+v", err)
				continue
			}
			log.Printf(" [x] Sent %s", data)
			return nil
		}
	}
}

func (p *Publisher) Connect() error {
	conn, err := amqp.Dial(p.address)
	if err != nil {
		return err
	}
	p.conn = conn

	p.channel, err = p.conn.Channel()
	if err != nil {
		return err
	}

	go func() {
		log.Printf("closing: %s", <-p.conn.NotifyClose(make(chan *amqp.Error)))
		p.done <- errors.New("channel closed")
	}()
	err = p.channel.QueueBind(
		p.queue,
		"",
		"",
		false,
		nil,
	)
	return err
}

func (p *Publisher) Reconnect() {
}

func (p *Publisher) Close() error {
	err := p.channel.Close()
	if err != nil {
		return err
	}
	return p.conn.Close()
}
