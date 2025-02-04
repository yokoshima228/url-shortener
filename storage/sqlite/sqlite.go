package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/yokoshima228/url-shortener/storage"
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

	stmt, err := s.db.Prepare("INSERT INTO url(alias, url) VALUES ($1, $2)")

	if err != nil {
		return fmt.Errorf("%s: %w", position, err)
	}

	_, err = stmt.Exec(alias, urlToSave)
	if err != nil {
		return fmt.Errorf("%s: %w", position, err)
	}

	return nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const position = "storage.sqlite.GetURL"

	stmt, err := s.db.Prepare(`SELECT url.url FROM url WHERE alias = $1`)
	if err != nil {
		return "", fmt.Errorf("%s: %w", position, err)
	}

	var url string
	err = stmt.QueryRow(alias).Scan(&url)
	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrUrlNotFound
	}
	if err != nil {
		return "", fmt.Errorf("%s: %w", position, err)
	}

	return url, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const position = "storage.sqlite.DeleteURL"

	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias = $1")
	if err != nil {
		return fmt.Errorf("%s: %w", position, err)
	}

	_, err = stmt.Exec(alias)
	return err
}
