package db

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Event struct {
	gorm.Model
	UUID        string
	Title       string
	Start       time.Time
	End         time.Time
	Description string
	OwnerID     string
	NotifyIn    time.Duration
}
