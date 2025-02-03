package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/mattn/go-sqlite3"
	"github.com/yokoshima228/url-shortener/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	const position = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", position, err)
	}
	defer db.Close()

	query, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", position, err)
	}

	_, err = query.Exec()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", position, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) error {
	const position = "storage.sqlite.SaveURL"

	_, err := s.db.Exec("INSERT INTO url(alias, url) VALUES ($1, $2)", alias, urlToSave)
	if err != nil {
		if err == sqlite3.ErrConstraintUnique {
			return storage.ErrUrlExists
		}

		return fmt.Errorf("%s: %w", position, err)
	}

	return nil
}
