package db

//nolint:golint
import (
	"time"

	"github.com/araddon/dateparse"
	in "github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal"
	cl "github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/app/calendar"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SQLDb struct {
	base *sqlx.DB
}

func NewSQLDatabase(cf config.DBConfiguration) (cl.StorageInterface, error) {
	var err error
	var DB SQLDb
	DB.base, err = sqlx.Connect("mysql", cf.User+":"+cf.Password+"@("+cf.Host+":"+cf.Port+")/"+cf.Name+"?charset=utf8&parseTime=True&loc=Local")
	return &DB, err
}

func (db *SQLDb) validateTime(start time.Time, end time.Time, uuidExcept string) error {
	switch {
	case end.Before(start):
		return ErrEndBeforeStart
	case start.Before(time.Now()):
		return ErrEventInPast
	case end.Sub(start) < 5*time.Minute:
		return ErrTooShortEvent
	}

	ev := in.Event{}
	//nolint:sqlclosecheck
	rows, err := db.base.Queryx("SELECT * FROM events")
	if err != nil {
		return err
	}
	for rows.Next() {
		err := rows.StructScan(&ev)
		if err != nil {
			return err
		}
		s, err := dateparse.ParseAny(ev.Start)
		if err != nil {
			return err
		}
		e, err := dateparse.ParseAny(ev.End)
		if err != nil {
			return err
		}
		if !((s.Before(end) && e.Before(start)) || (start.Before(e) && end.Before(s))) {
			if ev.UUID != uuidExcept {
				return ErrDateBusy
			}
		}
	}
	return nil
}

func (db *SQLDb) CreateEvent(ev in.Event) (uuid.UUID, error) {
	id := uuid.New()
	start, err := dateparse.ParseAny(ev.Start)
	if err != nil {
		return id, err
	}
	end, err := dateparse.ParseAny(ev.End)
	if err != nil {
		return id, err
	}
	if err := db.validateTime(start, end, ""); err != nil {
		return id, err
	}
	ev.UUID = id.String()
	_, err = db.base.Exec(`INSERT INTO events (uuid, title, start, end, description, ownerid, notifyin) VALUES (?, ?, ?, ?, ?, ?, ?)`, id.String(), ev.Title, ev.Start, ev.End, ev.Description, ev.OwnerID, ev.NotifyIn)
	return id, err
}

func (db *SQLDb) GetAllEvents() ([]in.Event, error) {
	evs := []in.Event{}
	err := db.base.Select(&evs, "SELECT * FROM events")
	return evs, err
}

func (db *SQLDb) GetEventByUUID(id uuid.UUID) (in.Event, error) {
	ev := in.Event{}
	err := db.base.Get(&ev, "SELECT * FROM events WHERE uuid = ?", id.String())
	return ev, err
}

func (db *SQLDb) GetFromInterval(start time.Time, delta time.Duration) ([]in.Event, error) {
	res := []in.Event{}
	end := start.Add(delta)
	evs, err := db.GetAllEvents()
	if err != nil {
		return nil, err
	}
	for _, ev := range evs {
		s, err := dateparse.ParseAny(ev.Start)
		if err != nil {
			return nil, err
		}
		e, err := dateparse.ParseAny(ev.End)
		if err != nil {
			return nil, err
		}
		if (start.Before(s) || start == s) && (e.Before(end) || end == e) {
			res = append(res, ev)
		}
	}
	if len(res) != 0 {
		return res, nil
	}
	return nil, ErrEventNotFound
}

func (db *SQLDb) UpdateEvent(ev in.Event, id uuid.UUID) error {
	if _, err := db.GetEventByUUID(id); err != nil {
		return ErrEventNotFound
	}
	start, err := dateparse.ParseAny(ev.Start)
	if err != nil {
		return err
	}
	end, err := dateparse.ParseAny(ev.End)
	if err != nil {
		return err
	}
	if err := db.validateTime(start, end, id.String()); err != nil {
		return err
	}
	ev.UUID = id.String()
	_, err = db.base.NamedExec(`UPDATE events SET title=:title, start=:start, end=:end, description=:description, ownerid=:ownerid, notifyin=:notifyin WHERE :uuid = :uuid`, ev)
	return err
}

func (db *SQLDb) DeleteEvent(id uuid.UUID) error {
	if _, err := db.GetEventByUUID(id); err != nil {
		return ErrEventNotFound
	}
	_, err := db.base.Exec("DELETE FROM events WHERE uuid = ?", id)
	return err
}

func (db *SQLDb) DeleteAll() {
	_, _ = db.base.Exec("DELETE FROM events")
}
