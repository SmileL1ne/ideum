package sesm

import (
	"context"
	"forum/pkg/sesm/sqlit3store"
	"log"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
)

type SessionManager struct {
	DB         *sqlit3store.SQLite3Store
	Lifetime   time.Duration
	CookieName string
	ContextKey string
}

func New() *SessionManager {
	return &SessionManager{
		Lifetime:   12 * time.Hour,
		ContextKey: generateContextKey(),
		CookieName: "session",
	}
}

type sessionData struct {
	sessionID  string
	expiryTime time.Time
}

func (sm *SessionManager) LoadAndSave(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Cookie")

		var sessionID string
		cookie, err := r.Cookie(sm.CookieName)
		if err == nil {
			sessionID = cookie.Name
		}

		ctx, err := sm.Load(r.Context(), sessionID)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		sessionReq := r.WithContext(ctx)

		sessionWriter := &sessionWriter{
			ResponseWriter: w,
			request:        sessionReq,
			sessionManager: sm,
		}

		next.ServeHTTP(sessionWriter, sessionReq)

		if !sessionWriter.written {
			sm.WriteSessionCookie(w, sessionReq)
		}
	})
}

func (sm *SessionManager) Load(ctx context.Context, sessionID string) (context.Context, error) {
	if _, ok := ctx.Value(sm.ContextKey).(*sessionData); ok {
		return ctx, nil
	}

	if sessionID == "" {
		newSD, err := newSessionData(sm.Lifetime)
		if err != nil {
			return nil, err
		}
		return context.WithValue(ctx, sm.ContextKey, newSD), nil
	}

	expiryTime, found, err := sm.findInStore(ctx, sessionID)
	if err != nil {
		return nil, err
	} else if !found {
		newSD, err := newSessionData(sm.Lifetime)
		if err != nil {
			return nil, err
		}
		return context.WithValue(ctx, sm.ContextKey, newSD), nil
	}

	sd := &sessionData{expiryTime: expiryTime}

	return context.WithValue(ctx, sm.ContextKey, sd), nil
}

func createSessionID() (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

func newSessionData(lifetime time.Duration) (*sessionData, error) {
	newSessionID, err := createSessionID()
	if err != nil {
		return nil, err
	}
	return &sessionData{
		sessionID:  newSessionID,
		expiryTime: time.Now().Add(lifetime).UTC(),
	}, nil
}

type sessionWriter struct {
	http.ResponseWriter
	request        *http.Request
	sessionManager *SessionManager
	written        bool
}

func (sw *sessionWriter) Write(b []byte) (int, error) {
	if !sw.written {
		sw.sessionManager.WriteSessionCookie(sw.ResponseWriter, sw.request)
		sw.written = true
	}

	return sw.ResponseWriter.Write(b)
}

func (sw *sessionWriter) WriteHeader(code int) {
	if !sw.written {
		sw.sessionManager.WriteSessionCookie(sw.ResponseWriter, sw.request)
		sw.written = true
	}

	sw.ResponseWriter.WriteHeader(code)
}
