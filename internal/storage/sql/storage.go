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

type WordInfo struct {
	Info []byte `db:"documents_list"`
}

type LastInsertId struct {
	Id int
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

	rows, err := s.db.Queryx(`INSERT INTO document_source (url, title, data) VALUES ($1, $2, $3) RETURNING id`,
		document.Url, document.Title, document.Data)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	rows.Next()
	lastInsertId := &LastInsertId{}
	err = rows.StructScan(&lastInsertId)
	return lastInsertId.Id, err
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
	rows, err := s.db.Queryx(`SELECT documents_list FROM document_invert_index WHERE word = $1`, word)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, nil
	}
	wordInfo := &WordInfo{}
	err = rows.StructScan(&wordInfo)
	return wordInfo.Info, err
}

func (s *Storage) UpdateWordInfo(word string, wordInfo []byte) error {
	rows, err := s.db.Queryx(`INSERT INTO document_invert_index (word, documents_list) VALUES ($1, $2) 
	ON CONFLICT (word) DO UPDATE SET documents_list = excluded.documents_list`, word, wordInfo)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}
