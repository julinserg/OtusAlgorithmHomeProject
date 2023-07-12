package sqlstorage

import (
	"context"
	"fmt"

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

/*func (s *Storage) GetEventsByDay(date time.Time) ([]storage.Event, error) {
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
}*/

func (s *Storage) Add(document storage.DocumentSource) error {
	_, err := s.db.NamedExec(`INSERT INTO document_source (url)
	 VALUES (:url)`,
		map[string]interface{}{
			"url": document.Url,
		})
	return err
}

func (s *Storage) GetAllDocumentSource() ([]storage.DocumentSource, error) {
	docList := make([]storage.DocumentSource, 0)
	doc := storage.DocumentSource{}
	rows, err := s.db.Queryx(`SELECT * FROM document_source`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.StructScan(&doc)
		if err != nil {
			return nil, err
		}
		docList = append(docList, doc)
	}
	return docList, nil
}
