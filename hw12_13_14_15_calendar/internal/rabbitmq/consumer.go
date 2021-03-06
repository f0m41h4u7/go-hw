package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/cenkalti/backoff/v3"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/app/sender"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/streadway/amqp"
)

type Consumer struct {
	address string
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
	done    chan error
}

func NewConsumer() sender.ConsumerInterface {
	return &Consumer{
		address: "amqp://" + net.JoinHostPort(config.SendConf.Rabbit.Host, config.SendConf.Rabbit.Port),
		done:    make(chan error),
	}
}

func (c *Consumer) Connect() error {
	conn, err := amqp.Dial(c.address)
	if err != nil {
		return err
	}
	c.conn = conn

	c.channel, err = c.conn.Channel()
	if err != nil {
		return err
	}

	go func() {
		log.Printf("closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
		c.done <- errors.New("channel closed")
	}()
	log.Printf("connected to rabbit")

	q, err := c.channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	c.queue = q.Name

	if err = c.channel.ExchangeDeclare(
		"event",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	return c.channel.QueueBind(
		c.queue,
		"",
		"event",
		false,
		nil,
	)
}

func (c *Consumer) Reconnect(ctx context.Context) (<-chan amqp.Delivery, error) {
	be := backoff.NewExponentialBackOff()
	be.MaxElapsedTime = 3 * time.Minute
	be.InitialInterval = 1 * time.Second
	be.Multiplier = 2
	be.MaxInterval = 30 * time.Second

	b := backoff.WithContext(be, ctx)
	for {
		d := b.NextBackOff()
		if d == backoff.Stop {
			return nil, fmt.Errorf("stop reconnecting")
		}

		select {
		case <-ctx.Done():
			return nil, nil
		case <-time.After(d):
			if err := c.Connect(); err != nil {
				log.Printf("could not connect in reconnect call: %+v", err)

				continue
			}
			msgs, err := c.channel.Consume(
				c.queue,
				"",
				false,
				false,
				false,
				false,
				nil,
			)
			if err != nil {
				log.Printf("could not connect: %+v", err)

				continue
			}

			return msgs, nil
		}
	}
}

func (c *Consumer) Receive(ctx context.Context, fn func(<-chan amqp.Delivery)) error {
	msgs, err := c.Reconnect(ctx)
	if err != nil {
		return err
	}

	for {
		go fn(msgs)

		if <-c.done != nil {
			msgs, err = c.Reconnect(ctx)
			if err != nil {
				return err
			}
		}
	}
}

func (c *Consumer) Close() error {
	err := c.channel.Close()
	if err != nil {
		return err
	}

	return c.conn.Close()
}
