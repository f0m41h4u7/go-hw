package calendar

import (
	"fmt"
	"time"

	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal"
	"github.com/google/uuid"
)

const (
	day   = 24 * time.Hour
	week  = 7 * day
	month = 4 * week
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

func (c *Calendar) GetEventsForDay(date time.Time) ([]internal.Event, error) {
	return c.Storage.GetFromInterval(date, day)
}

func (c *Calendar) GetEventsForWeek(date time.Time) ([]internal.Event, error) {
	return c.Storage.GetFromInterval(date, week)
}

func (c *Calendar) GetEventsForMonth(date time.Time) ([]internal.Event, error) {
	return c.Storage.GetFromInterval(date, month)
}

func (c *Calendar) UpdateEvent(ev internal.Event, u uuid.UUID) error {
	return c.Storage.UpdateEvent(ev, u)
}

func (c *Calendar) DeleteEvent(u uuid.UUID) error {
	return c.Storage.DeleteEvent(u)
}
