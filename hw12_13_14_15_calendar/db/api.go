package db

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

var (
	ErrDateBusy      = errors.New("date is busy")
	ErrEventNotFound = errors.New("event not found")
)

type Event struct {
	gorm.Model
	Uuid        string
	Title       string
	Start       time.Time
	End         time.Time
	Description string
	OwnerId     string
	NotifyIn    time.Duration
}

func checkTime(start time.Time, end time.Time) bool {
	evs := []Event{}
	GetDB().Where("Start >= ? AND End <= ?", start, end).Find(&evs)
	if len(evs) == 0 {
		return true
	}
	return false
}

func CreateEvent(ev Event) error {
	if !checkTime(ev.Start, ev.End) {
		return ErrDateBusy
	}
	err := GetDB().Create(&ev).Error
	return err
}

func GetFromInterval(date time.Time, delta time.Duration) ([]Event, error) {
	evs := []Event{}
	GetDB().Where("Start >= ? AND End <= ?", date, date.Add(delta)).Find(&evs)
	if len(evs) == 0 {
		return nil, ErrEventNotFound
	}
	return evs, nil
}

func UpdateEvent(ev Event, id string) error {
	var event Event
	err := GetDB().Where("Uuid = ?", id).First(&event).Error
	if err != nil {
		return ErrEventNotFound
	}

	if !checkTime(ev.Start, ev.End) {
		return ErrDateBusy
	}

	event.Title = ev.Title
	event.Start = ev.Start
	event.End = ev.End
	event.Description = ev.Description
	event.OwnerId = ev.OwnerId
	event.NotifyIn = ev.NotifyIn

	GetDB().Save(&event)
	return nil
}

func DeleteEvent(id string) error {
	var event Event
	err := GetDB().Where("ID = ?", id).First(&event).Error
	if err != nil {
		return ErrEventNotFound
	}

	err = GetDB().Delete(event).Error
	if err != nil {
		return err
	}
	return nil
}
