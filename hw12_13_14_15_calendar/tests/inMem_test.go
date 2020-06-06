package tests

import (
	"testing"
	"time"

	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateInDB(t *testing.T) {
	base, _ := db.NewInMemDatabase()
	tmp, _ := time.ParseDuration("1h")
	t.Run("simple", func(t *testing.T) {
		UUID := uuid.New()
		ev := db.Event{
			UUID:        UUID.String(),
			Title:       "cool event",
			Start:       time.Now().Add(1 * time.Hour),
			End:         time.Now().Add(6 * time.Hour),
			Description: "test",
			OwnerID:     uuid.New().String(),
			NotifyIn:    tmp,
		}

		err := base.CreateEvent(ev)
		require.Nil(t, err)

		event, err := base.GetEventByUUID(UUID)
		require.Nil(t, err)
		require.Equal(t, ev, event)
	})

	t.Run("date busy", func(t *testing.T) {
		ev := db.Event{
			UUID:        uuid.New().String(),
			Title:       "cool event",
			Start:       time.Now().Add(2 * time.Hour),
			End:         time.Now().Add(10 * time.Hour),
			Description: "test",
			OwnerID:     uuid.New().String(),
			NotifyIn:    tmp,
		}

		err := base.CreateEvent(ev)
		require.Equal(t, db.ErrDateBusy, err)
	})

	t.Run("end before start", func(t *testing.T) {
		ev := db.Event{
			UUID:        uuid.New().String(),
			Title:       "cool event",
			Start:       time.Now().Add(15 * time.Hour),
			End:         time.Now().Add(1 * time.Hour),
			Description: "test",
			OwnerID:     uuid.New().String(),
			NotifyIn:    tmp,
		}

		err := base.CreateEvent(ev)
		require.Equal(t, db.ErrEndBeforeStart, err)
	})

	t.Run("event in the past", func(t *testing.T) {
		ev := db.Event{
			UUID:        uuid.New().String(),
			Title:       "best event 2",
			Start:       time.Now().AddDate(-1, 0, 0),
			End:         time.Now(),
			Description: "test",
			OwnerID:     uuid.New().String(),
			NotifyIn:    tmp,
		}

		err := base.CreateEvent(ev)
		require.Equal(t, db.ErrEventInPast, err)
	})

	t.Run("end before start", func(t *testing.T) {
		ev := db.Event{
			UUID:        uuid.New().String(),
			Title:       "cool event",
			Start:       time.Now().AddDate(1, 0, 0),
			End:         time.Now().AddDate(1, 0, 0).Add(2 * time.Minute),
			Description: "test",
			OwnerID:     uuid.New().String(),
			NotifyIn:    tmp,
		}

		err := base.CreateEvent(ev)
		require.Equal(t, db.ErrTooShortEvent, err)
	})
}

func TestGetFromDB(t *testing.T) {
	base, _ := db.NewInMemDatabase()
	tmp, _ := time.ParseDuration("1h")
	t.Run("get from interval", func(t *testing.T) {
		date := time.Date(2023, 3, 11, 9, 0, 0, 0, time.UTC)
		ev1 := db.Event{
			UUID:        uuid.New().String(),
			Title:       "cool event",
			Start:       date,
			End:         date.Add(3 * time.Hour),
			Description: "test",
			OwnerID:     uuid.New().String(),
			NotifyIn:    tmp,
		}

		ev2 := db.Event{
			UUID:        uuid.New().String(),
			Title:       "next cool event",
			Start:       date.Add(3 * 24 * time.Hour),
			End:         date.Add(3 * 24 * time.Hour).Add(5 * time.Hour),
			Description: "test",
			OwnerID:     uuid.New().String(),
			NotifyIn:    tmp,
		}

		err := base.CreateEvent(ev1)
		require.Nil(t, err)
		err = base.CreateEvent(ev2)
		require.Nil(t, err)

		evs, err := base.GetFromInterval(date, 10*24*time.Hour)
		require.Nil(t, err)
		require.Equal(t, 2, len(evs))

		_, err = base.GetFromInterval(time.Now(), 3*24*time.Hour)
		require.Equal(t, db.ErrEventNotFound, err)
	})
}

func TestUpdateEvent(t *testing.T) {
	base, _ := db.NewInMemDatabase()
	tmp, _ := time.ParseDuration("1h")
	t.Run("update existing event", func(t *testing.T) {
		date := time.Date(2023, 3, 11, 9, 0, 0, 0, time.UTC)
		EventUUID := uuid.New()
		ev1 := db.Event{
			UUID:        EventUUID.String(),
			Title:       "cool event",
			Start:       date,
			End:         date.Add(3 * time.Hour),
			Description: "test",
			OwnerID:     uuid.New().String(),
			NotifyIn:    tmp,
		}
		_ = base.CreateEvent(ev1)

		ev2 := db.Event{
			UUID:        EventUUID.String(),
			Title:       "new event",
			Start:       date,
			End:         date.Add(3 * time.Hour),
			Description: "test",
			OwnerID:     uuid.New().String(),
			NotifyIn:    tmp,
		}

		err := base.UpdateEvent(ev2, EventUUID)
		require.Nil(t, err)
		event, _ := base.GetEventByUUID(EventUUID)
		require.Equal(t, ev2, event)
	})

	t.Run("update nonexistent event", func(t *testing.T) {
		EventUUID := uuid.New()
		ev := db.Event{
			UUID:        EventUUID.String(),
			Title:       "some new event",
			Start:       time.Now(),
			End:         time.Now().Add(3 * time.Hour),
			Description: "test",
			OwnerID:     uuid.New().String(),
			NotifyIn:    tmp,
		}

		err := base.UpdateEvent(ev, EventUUID)
		require.Equal(t, db.ErrEventNotFound, err)
	})
}

func TestDeleteEvent(t *testing.T) {
	base, _ := db.NewInMemDatabase()
	tmp, _ := time.ParseDuration("1h")
	t.Run("delete existing event", func(t *testing.T) {
		date := time.Date(2023, 3, 11, 9, 0, 0, 0, time.UTC)
		EventUUID := uuid.New()
		ev1 := db.Event{
			UUID:        EventUUID.String(),
			Title:       "cool event",
			Start:       date,
			End:         date.Add(3 * time.Hour),
			Description: "test",
			OwnerID:     uuid.New().String(),
			NotifyIn:    tmp,
		}

		ev2 := db.Event{
			UUID:        uuid.New().String(),
			Title:       "next cool event",
			Start:       date.Add(3 * 24 * time.Hour),
			End:         date.Add(3 * 24 * time.Hour).Add(5 * time.Hour),
			Description: "test",
			OwnerID:     uuid.New().String(),
			NotifyIn:    tmp,
		}

		_ = base.CreateEvent(ev1)
		_ = base.CreateEvent(ev2)

		err := base.DeleteEvent(EventUUID)
		require.Nil(t, err)
		evs, _ := base.GetAllEvents()
		require.Equal(t, 1, len(evs))
	})

	t.Run("delete nonexistent event", func(t *testing.T) {
		EventUUID := uuid.New()

		err := base.DeleteEvent(EventUUID)
		require.Equal(t, db.ErrEventNotFound, err)
	})
}
