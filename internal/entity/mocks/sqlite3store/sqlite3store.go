package sqlite3store

import (
	"context"
	"time"
)

var (
	userIDMock = 1
	expiryMock = time.Now().Add(1_000_000 * time.Hour)
)

type SQLite3StoreMock struct {
}

func New() *SQLite3StoreMock {
	return &SQLite3StoreMock{}
}

func (s *SQLite3StoreMock) StoreFind(ctx context.Context, sessionID string) (int, time.Time, error) {
	return userIDMock, expiryMock, nil
}

func (s *SQLite3StoreMock) StoreCommit(ctx context.Context, sessionID string, userID int, expiry time.Time) error {
	return nil
}

func (s *SQLite3StoreMock) StoreDeleteAll(ctx context.Context, userID int) error {
	return nil
}

func (s *SQLite3StoreMock) StoreDelete(ctx context.Context, sessionID string) error {
	return nil
}
