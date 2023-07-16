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

func (s *Storage) Add(document storage.Document) (int, error) {
	var lastInsertId int
	err := s.db.QueryRowx(`INSERT INTO document_source (url, title, data) VALUES ($1, $2, $3) RETURNING id`,
		document.Url, document.Title, document.Data).Scan(&lastInsertId)
	if err != nil {
		return 0, err
	}
	return lastInsertId, err
}

func (s *Storage) GetAllDocumentSource() ([]storage.Document, error) {
	docList := make([]storage.Document, 0)
	doc := storage.Document{}
	rows, err := s.db.Queryx(`SELECT id, url, title FROM document_source ORDER BY id`)
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

func (s *Storage) GetWordInfo(word string) ([]byte, error) {
	return nil, nil
}

func (s *Storage) UpdateWordInfo(word string, wordInfo []byte) error {
	return nil
}
