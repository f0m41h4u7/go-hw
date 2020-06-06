package db

//nolint: golint
import (
	"time"

	cfg "github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/config"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type SQLDb struct {
	base *gorm.DB
}

func NewSQLDatabase() (Database, error) {
	var err error
	connectAddr := cfg.Conf.Database.User + ":" + cfg.Conf.Database.Password + "@(" + cfg.Conf.Database.Host + ":" + cfg.Conf.Database.Port + ")/default?charset=utf8&parseTime=True&loc=Local"

	var DB SQLDb
	DB.base, err = gorm.Open("mysql", connectAddr)
	if err != nil {
		return &DB, err
	}
	DB.base.AutoMigrate(&Event{})
	return &DB, nil
}

func (db *SQLDb) GetAllEvents() ([]Event, error) {
	evs := []Event{}
	err := db.base.Find(&evs).Error
	return evs, err
}

func (db *SQLDb) GetEventByUUID(uuid uuid.UUID) (Event, error) {
	var event Event
	err := db.base.Where("UUID = ?", uuid.String()).First(&event).Error
	if err != nil {
		return event, ErrEventNotFound
	}
	return event, nil
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

	evs, err := db.GetAllEvents()
	if err != nil || len(evs) == 0 {
		return ErrEventNotFound
	}
	for _, ev := range evs {
		if !((ev.Start.Before(end) && ev.End.Before(start)) || (start.Before(ev.End) && end.Before(ev.Start))) {
			if ev.UUID != uuidExcept {
				return ErrDateBusy
			}
		}
	}
	return nil
}

func (db *SQLDb) CreateEvent(ev Event) error {
	if err := db.validateTime(ev.Start, ev.End, ""); err != nil {
		return err
	}
	err := db.base.Create(&ev).Error
	return err
}

func (db *SQLDb) GetFromInterval(start time.Time, delta time.Duration) ([]Event, error) {
	evs := []Event{}
	end := start.Add(delta)
	err := db.base.Where("Start >= ? AND End <= ?", start, end).Find(&evs).Error
	if err != nil {
		return nil, err
	}
	if len(evs) == 0 {
		return nil, ErrEventNotFound
	}
	return evs, nil
}

func (db *SQLDb) UpdateEvent(ev Event, uuid uuid.UUID) error {
	var event Event
	err := db.base.Where("UUID = ?", uuid.String()).First(&event).Error
	if err != nil {
		return ErrEventNotFound
	}

	if err = db.validateTime(ev.Start, ev.End, uuid.String()); err != nil {
		return err
	}

	event.Title = ev.Title
	event.Start = ev.Start
	event.End = ev.End
	event.Description = ev.Description
	event.OwnerID = ev.OwnerID
	event.NotifyIn = ev.NotifyIn

	err = db.base.Save(&event).Error
	return err
}

func (db *SQLDb) DeleteEvent(uuid uuid.UUID) error {
	var event Event
	err := db.base.Where("UUID = ?", uuid.String()).First(&event).Error
	if err != nil {
		return ErrEventNotFound
	}

	err = db.base.Delete(event).Error
	if err != nil {
		return err
	}
	return nil
}
