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
		queue:   "eventQueue",
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
		c.done <- errors.New("channel Closed")
	}()

	err = c.channel.QueueBind(
		c.queue,
		"",
		"",
		false,
		nil,
	)
	return err
}

func (c *Consumer) Reconnect() (<-chan amqp.Delivery, error) {
	be := backoff.NewExponentialBackOff()
	be.MaxElapsedTime = time.Minute
	be.InitialInterval = 1 * time.Second
	be.Multiplier = 2
	be.MaxInterval = 15 * time.Second

	b := backoff.WithContext(be, context.Background())
	for {
		d := b.NextBackOff()
		if d == backoff.Stop {
			return nil, fmt.Errorf("stop reconnecting")
		}

		//nolint:gosimple
		select {
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
				fmt.Printf("could not connect: %+v", err)
				continue
			}

			return msgs, nil
		}
	}
}

func (c *Consumer) Receive(fn func(<-chan amqp.Delivery)) error {
	var err error
	if err = c.Connect(); err != nil {
		return err
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
		return err
	}

	for {
		go fn(msgs)

		if <-c.done != nil {
			msgs, err = c.Reconnect()
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
