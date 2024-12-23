package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to open the database: %w", op, err)
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS urls (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL
			);
		CREATE INDEX IF NOT EXISTS urls_alies_idx ON urls (alies);
		`)

	if err != nil {
		return nil, fmt.Errorf("%s: failed to prepare the statement: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to execute the statement: %w", op, err)
	}

	return &Storage{db: db}, nil
}
