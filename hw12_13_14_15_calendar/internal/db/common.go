package db

import "errors"

var (
	ErrDateBusy       = errors.New("date is busy")
	ErrEventNotFound  = errors.New("event not found")
	ErrEndBeforeStart = errors.New("end of event is before start")
	ErrEventInPast    = errors.New("cannot create event in past")
	ErrTooShortEvent  = errors.New("event should be at least 5 minutes long")
)
