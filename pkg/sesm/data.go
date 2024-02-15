package sesm

import (
	"context"
	"fmt"
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

	b, found, err := sm.Store.StoreFind(ctx, sessionID)
	if err != nil {
		return nil, err
	} else if !found {
		return context.WithValue(ctx, sm.ContextKey, newSessionData(sm.Lifetime)), nil
	}

	sd := &sessionData{
		sessionID: sessionID,
		status:    Unmodified,
	}
	if sd.expiryTime, sd.values, err = sm.decode(b); err != nil {
		return nil, err
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
	values     map[string]interface{}
	expiryTime time.Time
	mu         sync.Mutex
}

// newSessionData returs sessionData with given and default values
func newSessionData(lifetime time.Duration) *sessionData {
	return &sessionData{
		status:     Unmodified,
		values:     make(map[string]interface{}),
		expiryTime: time.Now().Add(lifetime).UTC().Add(6 * time.Hour),
	}
}

// RenewToken updates tocken for given context, deleting old one (if exists)
func (sm *SessionManager) RenewToken(ctx context.Context) error {
	sd := sm.getSessionDataFromContext(ctx)

	sd.mu.Lock()
	defer sd.mu.Unlock()

	if sd.sessionID != "" {
		err := sm.Store.StoreDelete(ctx, sd.sessionID)
		if err != nil {
			return err
		}
	}

	newSessionID, err := createSessionID()
	if err != nil {
		return err
	}

	sd.sessionID = newSessionID
	sd.expiryTime = time.Now().Add(sm.Lifetime).UTC().Add(6 * time.Hour)
	sd.status = Modified

	return nil
}

// Status returns current status of session data
func (sm *SessionManager) Status(ctx context.Context) Status {
	sd := sm.getSessionDataFromContext(ctx)

	sd.mu.Lock()
	defer sd.mu.Unlock()

	return sd.status
}

// Put puts value by key into 'values' field in session data.
//
// It sets session status to Modified.
func (sm *SessionManager) Put(ctx context.Context, key string, val interface{}) {
	sd := sm.getSessionDataFromContext(ctx)

	sd.mu.Lock()
	sd.values[key] = val
	sd.status = Modified
	sd.mu.Unlock()
}

// GetInt reads integer value from context by given key.
//
// If no value exists by this key or value is not int - it returns zero value
func (sm *SessionManager) GetInt(ctx context.Context, key string) int {
	val := sm.Get(ctx, key)
	num, ok := val.(int)
	if !ok {
		return 0
	}

	return num
}

// Get reads value from session data by given key.
//
// If no value exists, it returns nil
func (sm *SessionManager) Get(ctx context.Context, key string) interface{} {
	sd := sm.getSessionDataFromContext(ctx)

	// No need to use mutex here because the operation is read

	val, exists := sd.values[key]
	if !exists {
		return nil
	}

	return val
}

// PopString reads and deletes string value by given key from session data.
//
// If retrieved data is not string, it returns empty string.
//
// PopString sets status to Modified
func (sm *SessionManager) PopString(ctx context.Context, key string) string {
	val := sm.Pop(ctx, key)
	str, ok := val.(string)
	if !ok {
		return ""
	}
	return str
}

// Pop reads and deletes data by given key.
//
// If no value exists by given key, it returns nil.
//
// Pop sets status to Modified
func (sm *SessionManager) Pop(ctx context.Context, key string) interface{} {
	sd := sm.getSessionDataFromContext(ctx)

	sd.mu.Lock()
	defer sd.mu.Unlock()

	val, exists := sd.values[key]
	if !exists {
		return nil
	}
	delete(sd.values, key)
	sd.status = Modified

	return val
}

// Remove deletes key-value pair by given key
//
// Remove sets status to Modified
func (sm *SessionManager) Remove(ctx context.Context, key string) {
	sd := sm.getSessionDataFromContext(ctx)

	sd.mu.Lock()
	defer sd.mu.Unlock()

	delete(sd.values, key)
	sd.status = Modified
}

// Exists checks if value by given key exists in session data
func (sm *SessionManager) Exists(ctx context.Context, key string) bool {
	sd := sm.getSessionDataFromContext(ctx)
	_, ok := sd.values[key]
	return ok
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
