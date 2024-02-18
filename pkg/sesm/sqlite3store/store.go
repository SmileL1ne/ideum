package sqlite3store

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
)

var ErrNotFound = errors.New("sqlite3store: session not found")

type SQLite3Store struct {
	db *sql.DB
}

func New(db *sql.DB) *SQLite3Store {
	s := &SQLite3Store{db: db}
	go s.startCleanup(5 * time.Second)
	return s
}

func (s *SQLite3Store) StoreFind(ctx context.Context, sessionID string) (int, time.Time, error) {
	query := `
		SELECT user_id, expiry 
		FROM sessions 
		WHERE session_id = $1 AND datetime('now', 'localtime') < expiry
	`
	var userID int
	var expiry time.Time
	row := s.db.QueryRow(query, sessionID)

	err := row.Scan(&userID, &expiry)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, time.Time{}, ErrNotFound
		}
		return 0, time.Time{}, err
	}

	return userID, expiry, nil
}

func (s *SQLite3Store) StoreCommit(ctx context.Context, sessionID string, userID int, expiry time.Time) error {
	query := `
		REPLACE INTO sessions (session_id, user_id, expiry) 
		VALUES($1, $2, datetime($3))
	`

	formattedExpiry := expiry.Format("2006-01-02T15:04:05.999")

	_, err := s.db.Exec(query, sessionID, userID, formattedExpiry)
	if err != nil {
		return err
	}

	return nil
}

func (s *SQLite3Store) StoreDeleteAll(ctx context.Context, userID int) error {
	query := `
		DELETE FROM sessions WHERE user_id = $1
	`
	_, err := s.db.ExecContext(ctx, query, userID)
	return err
}

func (s *SQLite3Store) StoreDelete(ctx context.Context, sessionID string) error {
	query := `
		DELETE FROM sessions WHERE session_id = $1
	`
	_, err := s.db.ExecContext(ctx, query, sessionID)
	return err
}

func (s *SQLite3Store) startCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		err := s.deleteExpired()
		if err != nil {
			log.Println(err)
		}
	}
}

func (s *SQLite3Store) deleteExpired() error {
	_, err := s.db.Exec("DELETE FROM sessions WHERE expiry < datetime('now', 'localtime')")
	return err
}
