package sqlit3store

import (
	"database/sql"
	"time"
)

/*
	TODO:
	1. Add background cleanup of expired tokens
	reference - https://github.com/alexedwards/scs/blob/master/sqlite3store/sqlite3store.go
*/

type SQLite3Store struct {
	db *sql.DB
	// here add stopCleanup chan
}

func New(db *sql.DB) *SQLite3Store {
	return &SQLite3Store{db: db}
}

func (s *SQLite3Store) Find(sessionID string) (time.Time, bool, error) {

}
