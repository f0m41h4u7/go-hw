package scheduler

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/araddon/dateparse"
	in "github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/app/calendar"
	"github.com/google/uuid"
)

const year = 365 * 24 * time.Hour

type Scheduler struct {
	Storage   calendar.StorageInterface
	Publisher PublisherInterface
}

func NewScheduler(st calendar.StorageInterface, p PublisherInterface) Scheduler {
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
	log.Printf("Sending data %s\n", data)
	return s.Publisher.Send(data)
}

func (s *Scheduler) Scan() error {
	evs, err := s.Storage.GetAllEvents()
	if err != nil {
		return err
	}
	fmt.Println(evs)
	for _, ev := range evs {
		res, err := checkYear(ev.End)
		if err != nil {
			return err
		}
		if !res {
			id, err := uuid.Parse(ev.UUID)
			if err != nil {
				return err
			}
			log.Printf("Deleting old event %s\n", ev.UUID)
			_ = s.Storage.DeleteEvent(id)
		}

		if ev.NotifyIn != "" {
			start, err := dateparse.ParseAny(ev.Start)
			if err != nil {
				return err
			}
			notif, err := strconv.Atoi(ev.NotifyIn)
			if err != nil {
				return err
			}
			if time.Until(start) <= time.Duration(notif) {
				return s.Publish(ev)
			}
		}
	}

	return nil
}

func (s *Scheduler) Stop() error {
	return s.Publisher.Close()
}
