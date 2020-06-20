package db

import (
	"testing"
	"time"

	in "github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func Test_InMem_Create(t *testing.T) {
	base, _ := NewInMemDatabase()
	tmp, _ := time.ParseDuration("1h")
	t.Run("simple", func(t *testing.T) {
		ev := in.Event{
			Title:       "cool event",
			Start:       time.Now().Add(1 * time.Hour).String(),
			End:         time.Now().Add(6 * time.Hour).String(),
			Description: "test",
			OwnerID:     uuid.New().String(),
			NotifyIn:    tmp.String(),
		}

		id, err := base.CreateEvent(ev)
		require.Nil(t, err)
		ev.UUID = id.String()

		event, err := base.GetEventByUUID(id)
		require.Nil(t, err)
		require.Equal(t, ev, event)
	})

	t.Run("date busy", func(t *testing.T) {
		ev := in.Event{
			Title:       "cool event",
			Start:       time.Now().Add(2 * time.Hour).String(),
			End:         time.Now().Add(10 * time.Hour).String(),
			Description: "test",
			OwnerID:     uuid.New().String(),
			NotifyIn:    tmp.String(),
		}

		_, err := base.CreateEvent(ev)
		require.Equal(t, ErrDateBusy, err)
	})

	t.Run("end before start", func(t *testing.T) {
		ev := in.Event{
			Title:       "cool event",
			Start:       time.Now().Add(15 * time.Hour).String(),
			End:         time.Now().Add(1 * time.Hour).String(),
			Description: "test",
			OwnerID:     uuid.New().String(),
			NotifyIn:    tmp.String(),
		}

		_, err := base.CreateEvent(ev)
		require.Equal(t, ErrEndBeforeStart, err)
	})

	t.Run("event in the past", func(t *testing.T) {
		ev := in.Event{
			Title:       "cool event",
			Start:       time.Now().AddDate(-1, 0, 0).String(),
			End:         time.Now().String(),
			Description: "test",
			OwnerID:     uuid.New().String(),
			NotifyIn:    tmp.String(),
		}

		_, err := base.CreateEvent(ev)
		require.Equal(t, ErrEventInPast, err)
	})

	t.Run("end before start", func(t *testing.T) {
		ev := in.Event{
			Title:       "cool event",
			Start:       time.Now().AddDate(1, 0, 0).String(),
			End:         time.Now().AddDate(1, 0, 0).Add(2 * time.Minute).String(),
			Description: "test",
			OwnerID:     uuid.New().String(),
			NotifyIn:    tmp.String(),
		}

		_, err := base.CreateEvent(ev)
		require.Equal(t, ErrTooShortEvent, err)
	})
}

func Test_InMem_Get(t *testing.T) {
	base, _ := NewInMemDatabase()
	tmp, _ := time.ParseDuration("1h")
	t.Run("get from interval", func(t *testing.T) {
		date := time.Date(2023, 3, 11, 9, 0, 0, 0, time.UTC)
		ev1 := in.Event{
			Title:       "cool event",
			Start:       date.Add(time.Hour).String(),
			End:         date.Add(3 * time.Hour).String(),
			Description: "test",
			OwnerID:     uuid.New().String(),
			NotifyIn:    tmp.String(),
		}

		ev2 := in.Event{
			Title:       "next cool event",
			Start:       date.Add(3 * 24 * time.Hour).String(),
			End:         date.Add(3 * 24 * time.Hour).Add(5 * time.Hour).String(),
			Description: "test",
			OwnerID:     uuid.New().String(),
			NotifyIn:    tmp.String(),
		}

		_, err := base.CreateEvent(ev1)
		require.Nil(t, err)
		_, err = base.CreateEvent(ev2)
		require.Nil(t, err)

		evs, err := base.GetFromInterval(date, 10*24*time.Hour)
		require.Nil(t, err)
		require.Equal(t, 2, len(evs))

		_, err = base.GetFromInterval(time.Now(), 3*24*time.Hour)
		require.Equal(t, ErrEventNotFound, err)
	})
}

func Test_InMem_Update(t *testing.T) {
	base, _ := NewInMemDatabase()
	tmp, _ := time.ParseDuration("1h")
	t.Run("update existing event", func(t *testing.T) {
		date := time.Date(2023, 3, 11, 9, 0, 0, 0, time.UTC)
		ev1 := in.Event{
			Title:       "cool event",
			Start:       date.String(),
			End:         date.Add(3 * time.Hour).String(),
			Description: "test",
			OwnerID:     uuid.New().String(),
			NotifyIn:    tmp.String(),
		}
		id, _ := base.CreateEvent(ev1)

		ev2 := in.Event{
			Title:       "new event",
			Start:       date.String(),
			End:         date.Add(3 * time.Hour).String(),
			Description: "test",
			OwnerID:     uuid.New().String(),
			NotifyIn:    tmp.String(),
		}

		err := base.UpdateEvent(ev2, id)
		require.Nil(t, err)
		ev2.UUID = id.String()
		event, _ := base.GetEventByUUID(id)
		require.Equal(t, ev2, event)
	})

	t.Run("update nonexistent event", func(t *testing.T) {
		id := uuid.New()
		ev := in.Event{
			Title:       "some new event",
			Start:       time.Now().String(),
			End:         time.Now().Add(3 * time.Hour).String(),
			Description: "test",
			OwnerID:     uuid.New().String(),
			NotifyIn:    tmp.String(),
		}

		err := base.UpdateEvent(ev, id)
		require.Equal(t, ErrEventNotFound, err)
	})
}

func Test_InMem_Delete(t *testing.T) {
	base, _ := NewInMemDatabase()
	tmp, _ := time.ParseDuration("1h")
	t.Run("delete existing event", func(t *testing.T) {
		date := time.Date(2023, 3, 11, 9, 0, 0, 0, time.UTC)
		ev1 := in.Event{
			Title:       "cool event",
			Start:       date.String(),
			End:         date.Add(3 * time.Hour).String(),
			Description: "test",
			OwnerID:     uuid.New().String(),
			NotifyIn:    tmp.String(),
		}

		ev2 := in.Event{
			Title:       "next cool event",
			Start:       date.Add(3 * 24 * time.Hour).String(),
			End:         date.Add(3 * 24 * time.Hour).Add(5 * time.Hour).String(),
			Description: "test",
			OwnerID:     uuid.New().String(),
			NotifyIn:    tmp.String(),
		}

		id1, _ := base.CreateEvent(ev1)
		id2, _ := base.CreateEvent(ev2)
		ev2.UUID = id2.String()

		err := base.DeleteEvent(id1)
		require.Nil(t, err)
		evs, _ := base.GetAllEvents()
		require.Equal(t, 1, len(evs))
	})

	t.Run("delete nonexistent event", func(t *testing.T) {
		id := uuid.New()

		err := base.DeleteEvent(id)
		require.Equal(t, ErrEventNotFound, err)
	})
}
