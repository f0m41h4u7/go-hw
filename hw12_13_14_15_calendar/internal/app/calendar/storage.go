package calendar

import (
	"time"

	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal"
	"github.com/google/uuid"
)

type StorageInterface interface {
	CreateEvent(internal.Event) (uuid.UUID, error)
	GetAllEvents() ([]internal.Event, error)
	GetEventByUUID(uuid.UUID) (internal.Event, error)
	GetFromInterval(time.Time, time.Duration) ([]internal.Event, error)
	UpdateEvent(internal.Event, uuid.UUID) error
	DeleteEvent(uuid.UUID) error
	DeleteAll()
}
