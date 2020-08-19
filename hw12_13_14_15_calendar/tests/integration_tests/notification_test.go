package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/cucumber/messages-go/v10"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal"
	"github.com/streadway/amqp"
)

const amqpDSN = "amqp://guest:guest@rabbitmq:5672/"

var ev internal.Event

type notifyTest struct {
	conn          *amqp.Connection
	ch            *amqp.Channel
	messages      [][]byte
	messagesMutex *sync.RWMutex
	stopSignal    chan struct{}
	queue         string
}

func (test *notifyTest) startConsuming(*messages.Pickle) {
	test.messages = make([][]byte, 0)
	test.messagesMutex = new(sync.RWMutex)
	test.stopSignal = make(chan struct{})

	var err error
	ev = internal.Event{
		Title:       "Upcoming Event",
		Start:       time.Now().Add(20 * time.Minute).String(),
		End:         time.Now().Add(2 * time.Hour).String(),
		Description: "Some event that starts in 20 minutes and cannot be missed!",
		OwnerID:     "9bed7c53-c3bd-4f7e-92d1-5d98c04fb83a",
		NotifyIn:    "20m",
	}

	test.conn, err = amqp.Dial(amqpDSN)
	if err != nil {
		log.Fatal(err)
	}

	test.ch, err = test.conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	q, err := test.ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	test.queue = q.Name

	err = test.ch.ExchangeDeclare(
		"event",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	err = test.ch.QueueBind(
		test.queue,
		"",
		"event",
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	events, err := test.ch.Consume(
		test.queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	go func(stop <-chan struct{}) {
		for {
			select {
			case <-stop:
				return
			case event := <-events:
				test.messagesMutex.Lock()
				test.messages = append(test.messages, event.Body)
				test.messagesMutex.Unlock()
			}
		}
	}(test.stopSignal)
}

func (test *notifyTest) anEventIsAboutToStart() error {
	bytesEvent, err := ev.MarshalJSON()
	if err != nil {
		return err
	}
	r, err := http.Post("http://calendar:1337/create", "application/json", bytes.NewReader(bytesEvent))
	if err != nil {
		return err
	}
	if r.StatusCode != 200 {
		return fmt.Errorf("status code: %d", r.StatusCode)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	var resp getResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return err
	}
	ev.UUID = resp.UUID

	return nil
}

func (test *notifyTest) iReceiveNotification() error {
	time.Sleep(5 * time.Second)

	bytesEvent, err := ev.MarshalJSON()
	if err != nil {
		return err
	}
	str := string(bytesEvent)

	test.messagesMutex.RLock()
	defer test.messagesMutex.RUnlock()

	for _, msg := range test.messages {
		if string(msg) == str {
			return nil
		}
	}
	return fmt.Errorf("event with text '%s' was not found in %s", str, test.messages)
}

func (test *notifyTest) stopConsuming(*messages.Pickle, error) {
	close(test.stopSignal)

	err := test.ch.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = test.conn.Close()
	if err != nil {
		log.Fatal(err)
	}
	test.messages = nil
}
