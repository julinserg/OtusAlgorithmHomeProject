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

func (s *Storage) Add(document storage.DocumentSource) error {
	_, err := s.db.NamedExec(`INSERT INTO document_source (url, title, data)
	 VALUES (:url, :title, :data)`,
		map[string]interface{}{
			"url":   document.Url,
			"title": document.Title,
			"data":  document.Data,
		})
	return err
}

func (s *Storage) GetAllDocumentSource() ([]storage.DocumentSource, error) {
	docList := make([]storage.DocumentSource, 0)
	doc := storage.DocumentSource{}
	rows, err := s.db.Queryx(`SELECT id, url, title FROM document_source`)
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
