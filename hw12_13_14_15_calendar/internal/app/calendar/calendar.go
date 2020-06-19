package calendar

import (
	"fmt"
	"time"

	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal"
	"github.com/google/uuid"
)

type Calendar struct {
	Storage StorageInterface
}

func NewCalendar(st StorageInterface) Calendar {
	return Calendar{
		Storage: st,
	}
}

func (c *Calendar) CreateEvent(ev internal.Event) (uuid.UUID, error) {
	fmt.Println("create")
	return c.Storage.CreateEvent(ev)
}

func (c *Calendar) GetEventsFromDay(date time.Time) ([]internal.Event, error) {
	return c.Storage.GetFromInterval(date, 24*time.Hour)
}

func (c *Calendar) GetEventsFromWeek(date time.Time) ([]internal.Event, error) {
	return c.Storage.GetFromInterval(date, 7*24*time.Hour)
}

func (c *Calendar) GetEventsFromMonth(date time.Time) ([]internal.Event, error) {
	return c.Storage.GetFromInterval(date, 30*7*24*time.Hour)
}

func (c *Calendar) UpdateEvent(ev internal.Event, u uuid.UUID) error {
	return c.Storage.UpdateEvent(ev, u)
}

func (c *Calendar) DeleteEvent(u uuid.UUID) error {
	return c.Storage.DeleteEvent(u)
}
