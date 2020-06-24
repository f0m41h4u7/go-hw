package tests

import (
	"time"

	in "github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal"
	"github.com/google/uuid"
)

var TestEvent = in.Event{
	Title:       "cool event",
	Start:       time.Date(2023, 3, 11, 9, 0, 0, 0, time.UTC).String(),
	End:         time.Date(2023, 3, 11, 9, 0, 0, 0, time.UTC).Add(3 * time.Hour).String(),
	Description: "test",
	OwnerID:     uuid.New().String(),
	NotifyIn:    "1h",
}

type TestStorage struct {
	Err error
}

func (t *TestStorage) CreateEvent(ev in.Event) (uuid.UUID, error) {
	return uuid.New(), t.Err
}

func (t *TestStorage) GetAllEvents() ([]in.Event, error) {
	return nil, t.Err
}

func (t *TestStorage) GetEventByUUID(id uuid.UUID) (in.Event, error) {
	TestEvent.UUID = id.String()
	return TestEvent, t.Err
}

func (t *TestStorage) GetFromInterval(st time.Time, del time.Duration) ([]in.Event, error) {
	evs := []in.Event{}
	evs = append(evs, TestEvent)
	evs = append(evs, TestEvent)
	return evs, t.Err
}

func (t *TestStorage) UpdateEvent(ev in.Event, id uuid.UUID) error {
	return t.Err
}

func (t *TestStorage) DeleteEvent(id uuid.UUID) error {
	return t.Err
}

func (t *TestStorage) DeleteAll() {}
