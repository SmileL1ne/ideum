package sesm

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gofrs/uuid"
)

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
	if sd.expiryTime, sd.values, err = sm.Decode(b); err != nil {
		return nil, err
	}

	return context.WithValue(ctx, sm.ContextKey, sd), nil
}

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

func newSessionData(lifetime time.Duration) *sessionData {
	return &sessionData{
		status:     Unmodified,
		values:     make(map[string]interface{}),
		expiryTime: time.Now().Add(lifetime).UTC().Add(6 * time.Hour), // Change to UTC + 6
	}
}

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

func (sm *SessionManager) Status(ctx context.Context) Status {
	sd := sm.getSessionDataFromContext(ctx)

	sd.mu.Lock()
	defer sd.mu.Unlock()

	return sd.status
}

// Puts value by key to the 'values' field in session data. Sets session status to Modified
func (sm *SessionManager) Put(ctx context.Context, key string, val interface{}) {
	sd := sm.getSessionDataFromContext(ctx)

	sd.mu.Lock()
	sd.values[key] = val
	sd.status = Modified
	sd.mu.Unlock()
}

func (sm *SessionManager) PopString(ctx context.Context, key string) string {
	val := sm.Pop(ctx, key)
	str, ok := val.(string)
	if !ok {
		return ""
	}
	return str
}

func (sm *SessionManager) Remove(ctx context.Context, key string) {
	sd := sm.getSessionDataFromContext(ctx)

	sd.mu.Lock()
	defer sd.mu.Unlock()

	delete(sd.values, key)
	sd.status = Modified
}

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

func (sm *SessionManager) Exists(ctx context.Context, key string) bool {
	sd := sm.getSessionDataFromContext(ctx)
	_, ok := sd.values[key]
	return ok
}

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

func generateContextKey() contextKey {
	contextKeyMutex.Lock()
	defer contextKeyMutex.Unlock()
	atomic.AddUint64(&contextKeyID, 1)
	return contextKey(fmt.Sprintf("session.%d", contextKeyID))
}
