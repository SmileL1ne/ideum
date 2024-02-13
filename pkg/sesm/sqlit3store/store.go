package sqlit3store

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type SQLite3Store struct {
	db *sql.DB
}

func New(db *sql.DB) *SQLite3Store {
	s := &SQLite3Store{db: db}
	go s.startCleanup(5 * time.Second)
	return s
}

func (s *SQLite3Store) StoreFind(ctx context.Context, sessionID string) ([]byte, bool, error) {
	var b []byte
	row := s.db.QueryRow("SELECT data FROM sessions WHERE session_id = $1 AND datetime('now') < expiry", sessionID)

	err := row.Scan(&b)
	if err == sql.ErrNoRows {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}

	return b, true, nil
}

func (s *SQLite3Store) StoreCommit(ctx context.Context, sessionID string, b []byte, expiry time.Time) error {
	formattedExpiry := expiry.UTC().Format("2006-01-02T15:04:05.999")

	_, err := s.db.Exec("REPLACE INTO sessions (session_id, data, expiry) VALUES($1, $2, datetime($3))", sessionID, b, formattedExpiry)
	if err != nil {
		return err
	}

	return nil
}

func (s *SQLite3Store) StoreDelete(ctx context.Context, sessionID string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM sessions WHERE session_id = $1", sessionID)
	return err
}

func (s *SQLite3Store) startCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			err := s.deleteExpired()
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func (s *SQLite3Store) deleteExpired() error {
	_, err := s.db.Exec("DELETE FROM sessions WHERE expiry < datetime('now', 'localtime')")
	return err
}
