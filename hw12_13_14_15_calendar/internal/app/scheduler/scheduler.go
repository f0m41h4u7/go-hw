package scheduler

import (
	"fmt"
	"log"
	"time"

	"github.com/araddon/dateparse"
	in "github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal"
	"github.com/google/uuid"
)

const year = 365 * 24 * time.Hour

type Storage interface {
	GetAllEvents() ([]in.Event, error)
	DeleteEvent(uuid.UUID) error
}

type Scheduler struct {
	Storage   Storage
	Publisher PublisherInterface
}

func NewScheduler(st Storage, p PublisherInterface) Scheduler {
	return Scheduler{
		Storage:   st,
		Publisher: p,
	}
}

func checkYear(end string) (bool, error) {
	date, err := dateparse.ParseAny(end)
	if err != nil {
		return false, err
	}
	if time.Since(date) >= year {
		return false, nil
	}
	return true, nil
}

//nolint:interfacer
func (s *Scheduler) Publish(ev in.Event) error {
	data, err := ev.MarshalJSON()
	if err != nil {
		return err
	}
	log.Printf("sending notification %s\n", data)
	return s.Publisher.Send(data)
}

func (s *Scheduler) Scan() {
	evs, err := s.Storage.GetAllEvents()
	if err != nil {
		log.Printf("failed to get events: %+v\n", err)
	}
	fmt.Println(evs)
	for _, ev := range evs {
		res, err := checkYear(ev.End)
		if err != nil {
			log.Printf("error when checking event date: %+v\n", err)
		}
		if !res {
			id, err := uuid.Parse(ev.UUID)
			if err != nil {
				log.Printf("failed to parse uuid: %+v\n", err)
			}
			log.Printf("deleting old event %s\n", ev.UUID)
			_ = s.Storage.DeleteEvent(id)
		}
		//nolint:nestif
		if ev.NotifyIn != "" {
			start, err := dateparse.ParseAny(ev.Start)
			if err != nil {
				log.Printf("failed to parse date: %+v\n", err)
			}
			notif, err := time.ParseDuration(ev.NotifyIn)
			if err != nil {
				log.Printf("failed to parse notification interval: %+v\n", err)
			}
			if (time.Now().Before(start)) && (time.Until(start) <= notif) {
				err := s.Publish(ev)
				if err != nil {
					log.Printf("publisher error: %+v\n", err)
				}
			}
		}
	}
}

func (s *Scheduler) Stop() error {
	return s.Publisher.Close()
}
