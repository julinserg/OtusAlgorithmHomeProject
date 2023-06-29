package storage

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrEventIDNotSet       = errors.New("Event ID not set")
	ErrEventIDNotExist     = errors.New("Event ID not exist")
	ErrEventIDAlreadyExist = errors.New("Event ID already exist")
	ErrTimeBusy            = errors.New("time event is busy")
)

type Event struct {
	ID               string         `db:"id"`
	Title            string         `db:"title"`
	TimeStart        time.Time      `db:"time_start"`
	TimeEnd          time.Time      `db:"time_stop"`
	Description      sql.NullString `db:"description"`
	UserID           string         `db:"user_id"`
	NotificationTime sql.NullInt64  `db:"time_notify"`
}
