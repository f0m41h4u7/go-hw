package sender

import (
	"log"

	in "github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal"
	"github.com/streadway/amqp"
)

type Sender struct {
	Consumer ConsumerInterface
}

func NewSender(c ConsumerInterface) Sender {
	return Sender{
		Consumer: c,
	}
}

func getEvents(msgs <-chan amqp.Delivery) {
	for d := range msgs {
		ev := in.Event{}
		err := ev.UnmarshalJSON(d.Body)
		if err != nil {
			log.Printf("Cannot parse notification")
		} else {
			log.Printf("Received a notification: %s", d.Body)
		}
	}
}

func (s *Sender) Listen() error {
	return s.Consumer.Receive(getEvents)
}

func (s *Sender) Stop() error {
	return s.Consumer.Close()
}
