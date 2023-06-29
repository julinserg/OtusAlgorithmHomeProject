package sqlstorage

import (
	"context"
	"fmt"
	"strings"
	"time"

	// Register pgx driver for postgresql.
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/julinserg/OtusAlgorithmHomeProject/internal/storage"
)

type Storage struct {
	db *sqlx.DB
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context, dsn string) error {
	var err error
	s.db, err = sqlx.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("cannot open pgx driver: %w", err)
	}
	return s.db.PingContext(ctx)
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) get(id string) (storage.Event, error) {
	ev := storage.Event{}
	rows, err := s.db.NamedQuery(`SELECT id,title,time_start,time_stop,description,
	user_id,time_notify FROM events WHERE id=:id`, map[string]interface{}{"id": id})
	if err != nil {
		return ev, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.StructScan(&ev)
		if err != nil {
			return ev, err
		}
	}
	return ev, nil
}

func (s *Storage) GetEventsByDay(date time.Time) ([]storage.Event, error) {
	result := make([]storage.Event, 0)
	dateDayBegin := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	dateDayEnd := dateDayBegin.AddDate(0, 0, 1)
	rows, err := s.db.NamedQuery(`SELECT id,title,time_start,time_stop,description,
	user_id,time_notify FROM events WHERE time_start >= :timeS AND time_start < :timeE`,
		map[string]interface{}{
			"timeS": dateDayBegin,
			"timeE": dateDayEnd,
		})
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		ev := storage.Event{}
		err := rows.StructScan(&ev)
		if err != nil {
			return nil, err
		}
		result = append(result, ev)
	}
	return result, nil
}

func (s *Storage) getEventsByInterval(date1, date2 time.Time) ([]storage.Event, error) {
	result := make([]storage.Event, 0)
	rows, err := s.db.NamedQuery(`SELECT id,title,time_start,time_stop,description,
	user_id,time_notify FROM events WHERE time_start >= :timeS AND time_start <= :timeE`,
		map[string]interface{}{
			"timeS": date1,
			"timeE": date2,
		})
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		ev := storage.Event{}
		err := rows.StructScan(&ev)
		if err != nil {
			return nil, err
		}
		result = append(result, ev)
	}
	return result, nil
}

func (s *Storage) GetEventsByWeek(dateBeginWeek time.Time) ([]storage.Event, error) {
	dateEndWeek := dateBeginWeek.AddDate(0, 0, 7)
	return s.getEventsByInterval(dateBeginWeek, dateEndWeek)
}

func (s *Storage) GetEventsByMonth(dateBeginMonth time.Time) ([]storage.Event, error) {
	dateEndMonth := dateBeginMonth.AddDate(0, 1, 0)
	return s.getEventsByInterval(dateBeginMonth, dateEndMonth)
}

func (s *Storage) Add(event storage.Event) error {
	if len(event.ID) == 0 {
		return storage.ErrEventIDNotSet
	}
	_, err := s.db.NamedExec(`INSERT INTO events (id,title,time_start,time_stop,description,user_id,time_notify)
	 VALUES (:id,:title,:time_start,:time_stop,:description,:user_id,:time_notify)`,
		map[string]interface{}{
			"id":          event.ID,
			"title":       event.Title,
			"time_start":  event.TimeStart,
			"time_stop":   event.TimeEnd,
			"description": event.Description,
			"user_id":     event.UserID,
			"time_notify": event.NotificationTime,
		})
	if err != nil && strings.Contains(err.Error(), "time_start_unique") {
		return storage.ErrTimeBusy
	}
	return err
}

func (s *Storage) Update(event storage.Event) error {
	if len(event.ID) == 0 {
		return storage.ErrEventIDNotSet
	}
	result, err := s.db.NamedExec(`UPDATE events SET title=:title, time_start=:time_start,
	 time_stop=:time_stop,description=:description, 
	 user_id =:user_id, time_notify=:time_notify 
	 WHERE id = `+`'`+event.ID+`'`,
		map[string]interface{}{
			"title":       event.Title,
			"time_start":  event.TimeStart,
			"time_stop":   event.TimeEnd,
			"description": event.Description,
			"user_id":     event.UserID,
			"time_notify": event.NotificationTime,
		})
	if result != nil {
		rowAffected, errResult := result.RowsAffected()
		if err == nil && rowAffected == 0 && errResult == nil {
			return storage.ErrEventIDNotExist
		}
	}
	if err != nil && strings.Contains(err.Error(), "time_start_unique") {
		return storage.ErrTimeBusy
	}
	return err
}

func (s *Storage) Remove(id string) error {
	result, err := s.db.Exec(`DELETE FROM events	WHERE id = ` + `'` + id + `'`)
	rowAffected, errResult := result.RowsAffected()
	if err == nil && rowAffected == 0 && errResult == nil {
		return storage.ErrEventIDNotExist
	}
	return err
}

func (s *Storage) GetEventsForNotify(timeNow time.Time) ([]storage.Event, error) {
	// select id from events where time_start  <= to_timestamp('2022-10-23 01:05:00', 'YYYY-MM-DD HH24:MI:SS') + INTERVAL '1 min' * time_notify
	result := make([]storage.Event, 0)
	rows, err := s.db.Queryx(`SELECT id,title,time_start,time_stop,description,
	user_id,time_notify FROM events WHERE
	is_notifyed is NULL
	AND
	time_notify > 0
	AND 
	time_start <= to_timestamp('` + timeNow.String() + `', 'YYYY-MM-DD HH24:MI:SS') + (INTERVAL '1 milliseconds' * (time_notify/1000000))
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		ev := storage.Event{}
		err := rows.StructScan(&ev)
		if err != nil {
			return nil, err
		}
		result = append(result, ev)
	}
	return result, nil
}

func (s *Storage) MarkEventIsNotifyed(id string) error {
	result, err := s.db.Exec(`UPDATE events SET is_notifyed=TRUE WHERE id = ` + `'` + id + `'`)
	rowAffected, errResult := result.RowsAffected()
	if err == nil && rowAffected == 0 && errResult == nil {
		return storage.ErrEventIDNotExist
	}
	return err
}

func (s *Storage) RemoveOldYearEvent(timeLimit time.Time) (int64, error) {
	result, err := s.db.Exec(`DELETE FROM events WHERE id IN 
	(SELECT id FROM events WHERE time_start < to_timestamp('` + timeLimit.String() + `', 'YYYY-MM-DD HH24:MI:SS') - interval '1 year' LIMIT 100)`)
	rowAffected, errResult := result.RowsAffected()
	if err == nil {
		err = errResult
	}
	return rowAffected, err
}
