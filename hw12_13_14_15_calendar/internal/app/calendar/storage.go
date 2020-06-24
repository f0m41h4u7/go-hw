package calendar

import (
	"time"

	in "github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal"
	"github.com/google/uuid"
)

type StorageInterface interface {
	CreateEvent(in.Event) (uuid.UUID, error)
	GetAllEvents() ([]in.Event, error)
	GetEventByUUID(uuid.UUID) (in.Event, error)
	GetFromInterval(time.Time, time.Duration) ([]in.Event, error)
	UpdateEvent(in.Event, uuid.UUID) error
	DeleteEvent(uuid.UUID) error
	DeleteAll()
}
