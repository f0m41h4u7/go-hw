package db

import (
	"time"

	"github.com/google/uuid"
)

type InMemDB struct {
	base []Event
}

func NewInMemDatabase() (Database, error) {
	var im InMemDB
	im.base = []Event{}
	return &im, nil
}

func (im *InMemDB) GetAllEvents() ([]Event, error) {
	return im.base, nil
}

func (im *InMemDB) GetEventByUUID(uuid uuid.UUID) (Event, error) {
	if len(im.base) == 0 {
		return Event{}, ErrEventNotFound
	}
	for _, ev := range im.base {
		if ev.UUID == uuid.String() {
			return ev, nil
		}
	}
	return Event{}, ErrEventNotFound
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

	for _, ev := range im.base {
		if !((ev.Start.Before(end) && ev.End.Before(start)) || (start.Before(ev.End) && end.Before(ev.Start))) {
			if ev.UUID != uuidExcept {
				return ErrDateBusy
			}
		}
	}
	return nil
}

func (im *InMemDB) CreateEvent(ev Event) error {
	if err := im.validateTime(ev.Start, ev.End, ""); err != nil {
		return err
	}
	im.base = append(im.base, ev)
	return nil
}

func (im *InMemDB) GetFromInterval(start time.Time, delta time.Duration) ([]Event, error) {
	evs := []Event{}
	end := start.Add(delta)
	for _, ev := range im.base {
		if (start.Before(ev.Start) || start == ev.Start) && (ev.End.Before(end) || end == ev.End) {
			evs = append(evs, ev)
		}
	}
	if len(evs) != 0 {
		return evs, nil
	}
	return nil, ErrEventNotFound
}

func (im *InMemDB) UpdateEvent(newEvent Event, uuid uuid.UUID) error {
	for i, ev := range im.base {
		if ev.UUID == uuid.String() {
			im.base[i].Title = newEvent.Title
			im.base[i].Start = newEvent.Start
			im.base[i].End = newEvent.End
			im.base[i].Description = newEvent.Description
			im.base[i].OwnerID = newEvent.OwnerID
			im.base[i].NotifyIn = newEvent.NotifyIn
			if err := im.validateTime(newEvent.Start, newEvent.End, uuid.String()); err != nil {
				return err
			}
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
