package db

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrDateBusy       = errors.New("date is busy")
	ErrEventNotFound  = errors.New("event not found")
	ErrEndBeforeStart = errors.New("end of event is before start")
	ErrEventInPast    = errors.New("cannot create event in past")
	ErrTooShortEvent  = errors.New("event should be at least 5 minutes long")
)

type Database interface {
	CreateEvent(Event) error
	GetAllEvents() ([]Event, error)
	GetEventByUUID(uuid.UUID) (Event, error)
	GetFromInterval(time.Time, time.Duration) ([]Event, error)
	UpdateEvent(Event, uuid.UUID) error
	DeleteEvent(uuid.UUID) error
}
