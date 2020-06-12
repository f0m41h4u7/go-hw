package db

import (
	"time"

	"github.com/araddon/dateparse"
	in "github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal"
	cl "github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/app/calendar"
	"github.com/google/uuid"
)

type InMemDB struct {
	base []in.Event
}

func NewInMemDatabase() (cl.StorageInterface, error) {
	var im InMemDB
	im.base = []in.Event{}
	return &im, nil
}

func (im *InMemDB) GetAllEvents() ([]in.Event, error) {
	return im.base, nil
}

func (im *InMemDB) GetEventByUUID(uuid uuid.UUID) (in.Event, error) {
	if len(im.base) == 0 {
		return in.Event{}, ErrEventNotFound
	}
	for _, ev := range im.base {
		if ev.UUID == uuid.String() {
			return ev, nil
		}
	}
	return in.Event{}, ErrEventNotFound
}

func (im *InMemDB) validateTime(start time.Time, end time.Time, uuidExcept string) error {
	switch {
	case end.Before(start):
		return ErrEndBeforeStart
	case start.Before(time.Now()):
		return ErrEventInPast
	case end.Sub(start) < 5*time.Minute:
		return ErrTooShortEvent
	}

	var s, e time.Time
	for _, ev := range im.base {
		s, _ = dateparse.ParseAny(ev.Start)
		e, _ = dateparse.ParseAny(ev.End)
		if !((s.Before(end) && e.Before(start)) || (start.Before(e) && end.Before(s))) {
			if ev.UUID != uuidExcept {
				return ErrDateBusy
			}
		}
	}
	return nil
}

func (im *InMemDB) CreateEvent(ev in.Event) (uuid.UUID, error) {
	id := uuid.New()
	start, _ := dateparse.ParseAny(ev.Start)
	end, _ := dateparse.ParseAny(ev.End)
	if err := im.validateTime(start, end, ""); err != nil {
		return id, err
	}
	ev.UUID = id.String()
	im.base = append(im.base, ev)
	return id, nil
}

func (im *InMemDB) GetFromInterval(start time.Time, delta time.Duration) ([]in.Event, error) {
	evs := []in.Event{}
	end := start.Add(delta)
	var s, e time.Time
	for _, ev := range im.base {
		s, _ = dateparse.ParseAny(ev.Start)
		e, _ = dateparse.ParseAny(ev.End)
		if (start.Before(s) || start == s) && (e.Before(end) || end == e) {
			evs = append(evs, ev)
		}
	}
	if len(evs) != 0 {
		return evs, nil
	}
	return nil, ErrEventNotFound
}

func (im *InMemDB) UpdateEvent(newEvent in.Event, uuid uuid.UUID) error {
	var s, e time.Time
	for i, ev := range im.base {
		if ev.UUID == uuid.String() {
			s, _ = dateparse.ParseAny(newEvent.Start)
			e, _ = dateparse.ParseAny(newEvent.End)
			if err := im.validateTime(s, e, uuid.String()); err != nil {
				return err
			}

			im.base[i].Title = newEvent.Title
			im.base[i].Start = newEvent.Start
			im.base[i].End = newEvent.End
			im.base[i].Description = newEvent.Description
			im.base[i].OwnerID = newEvent.OwnerID
			im.base[i].NotifyIn = newEvent.NotifyIn
			return nil
		}
	}
	return ErrEventNotFound
}

func (im *InMemDB) DeleteEvent(uuid uuid.UUID) error {
	for i, ev := range im.base {
		if ev.UUID == uuid.String() {
			im.base[i] = im.base[len(im.base)-1]
			im.base = im.base[:len(im.base)-1]
			return nil
		}
	}
	return ErrEventNotFound
}

func (im *InMemDB) DeleteAll() {
	im.base = nil
}
