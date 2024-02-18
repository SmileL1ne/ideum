package sesm

import (
	"context"
	"errors"
	"fmt"
	"forum/pkg/sesm/sqlite3store"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gofrs/uuid"
)

// Load returns context with loaded session.
//
// If empty sessionID is given it returns context with new session created.
func (sm *SessionManager) Load(ctx context.Context, sessionID string) (context.Context, error) {
	if _, ok := ctx.Value(sm.ContextKey).(*sessionData); ok {
		return ctx, nil
	}

	if sessionID == "" {
		return context.WithValue(ctx, sm.ContextKey, newSessionData(sm.Lifetime)), nil
	}

	userID, expiry, err := sm.Store.StoreFind(ctx, sessionID)
	if err != nil {
		if errors.Is(err, sqlite3store.ErrNotFound) {
			sd := newSessionData(time.Second)
			sd.status = Destroyed
			return context.WithValue(ctx, sm.ContextKey, sd), nil
		}
		return nil, err
	}

	sd := &sessionData{
		sessionID:  sessionID,
		status:     Unmodified,
		userID:     userID,
		expiryTime: expiry,
	}

	return context.WithValue(ctx, sm.ContextKey, sd), nil
}

// createSessionID generates session id using uuid package
func createSessionID() (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

type Status int

const (
	Unmodified Status = iota

	Modified

	Destroyed
)

type sessionData struct {
	sessionID  string
	status     Status
	userID     int
	expiryTime time.Time
	mu         sync.Mutex
}

// newSessionData returs sessionData with given lifetime and default values
func newSessionData(lifetime time.Duration) *sessionData {
	return &sessionData{
		status:     Unmodified,
		expiryTime: time.Now().Local().Add(lifetime),
	}
}

// RenewToken updates token for given context, deleting old one (if exists)
func (sm *SessionManager) RenewToken(ctx context.Context, userID int) error {
	sd := sm.getSessionDataFromContext(ctx)

	sd.mu.Lock()
	defer sd.mu.Unlock()

	err := sm.Store.StoreDeleteAll(ctx, userID)
	if err != nil {
		return err
	}

	newSessionID, err := createSessionID()
	if err != nil {
		return err
	}

	sd.sessionID = newSessionID
	sd.expiryTime = time.Now().Add(sm.Lifetime)
	sd.status = Modified

	return nil
}

// DeleteToken deletes token from database
func (sm *SessionManager) DeleteToken(ctx context.Context) error {
	sd := sm.getSessionDataFromContext(ctx)

	sd.mu.Lock()
	defer sd.mu.Unlock()

	if sd.sessionID != "" {
		err := sm.Store.StoreDelete(ctx, sd.sessionID)
		if err != nil {
			return err
		}
	}

	sd.status = Destroyed

	return nil
}

// Status returns current status of session data
func (sm *SessionManager) Status(ctx context.Context) Status {
	sd := sm.getSessionDataFromContext(ctx)

	sd.mu.Lock()
	defer sd.mu.Unlock()

	return sd.status
}

// PutUserID puts userID in session data.
//
// It sets session status to Modified.
func (sm *SessionManager) PutUserID(ctx context.Context, userID int) {
	sd := sm.getSessionDataFromContext(ctx)

	sd.mu.Lock()
	sd.userID = userID
	sd.status = Modified
	sd.mu.Unlock()
}

// GetUserID reads userID from session data.
func (sm *SessionManager) GetUserID(ctx context.Context) int {
	sd := sm.getSessionDataFromContext(ctx)

	// No need to use mutex here because the operation is read

	return sd.userID
}

// Exists checks if userID exists in session data
func (sm *SessionManager) ExistsUserID(ctx context.Context) bool {
	sd := sm.getSessionDataFromContext(ctx)
	return sd.userID != 0
}

// getSessionDataFromContext retrieves session data from given contex.
//
// It panics if no session data is found
func (sm *SessionManager) getSessionDataFromContext(ctx context.Context) *sessionData {
	sd, ok := ctx.Value(sm.ContextKey).(*sessionData)
	if !ok {
		panic("sesm: no session data in context")
	}
	return sd
}

type contextKey string

var (
	contextKeyID    uint64
	contextKeyMutex = &sync.Mutex{}
)

// generateContextKey gerenerates new ContextKey
func generateContextKey() contextKey {
	contextKeyMutex.Lock()
	defer contextKeyMutex.Unlock()
	atomic.AddUint64(&contextKeyID, 1)
	return contextKey(fmt.Sprintf("session.%d", contextKeyID))
}
