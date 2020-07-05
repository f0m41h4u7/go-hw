package sender

import (
	"bytes"
	"log"
	"time"

	in "github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal"
)

type Sender struct {
	Consumer ConsumerInterface
}

func NewSender(c ConsumerInterface) Sender {
	return Sender{
		Consumer: c,
	}
}

func (s *Sender) Listen() error {
	msgs, err := s.Consumer.Receive()
	if err != nil {
		return err
	}
	for d := range msgs {
		ev := in.Event{}
		err := ev.UnmarshalJSON(d.Body)
		if err != nil {
			return err
		}
		log.Printf("Received a message: %s", d.Body)
		dotCount := bytes.Count(d.Body, []byte("."))
		t := time.Duration(dotCount)
		time.Sleep(t * time.Second)
		err = d.Ack(false)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Sender) Stop() error {
	return s.Consumer.Close()
}
