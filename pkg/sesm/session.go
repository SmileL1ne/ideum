package sesm

import (
	"context"
	"forum/pkg/sesm/sqlit3store"
	"log"
	"net/http"
	"time"
)

type SessionManager struct {
	Store      *sqlit3store.SQLite3Store
	Lifetime   time.Duration
	CookieName string
	ContextKey contextKey
}

func New() *SessionManager {
	return &SessionManager{
		Lifetime:   12 * time.Hour,
		ContextKey: generateContextKey(),
		CookieName: "session",
	}
}

func (sm *SessionManager) LoadAndSave(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Cookie")

		var sessionID string
		cookie, err := r.Cookie(sm.CookieName)
		if err == nil {
			sessionID = cookie.Value
		}

		ctx, err := sm.Load(r.Context(), sessionID)
		if err != nil {
			log.Println("LoadAndSave:", err.Error())
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
			sm.commitAndWriteSessionCookie(sessionWriter, sessionReq)
		}
	})
}

func (sm *SessionManager) commitAndWriteSessionCookie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch sm.Status(ctx) {
	case Modified:
		token, expiry, err := sm.Commit(ctx)
		if err != nil {
			log.Println("Commit session:", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		sm.WriteSessionCookie(ctx, w, token, expiry)
	}
}

func (sm *SessionManager) Commit(ctx context.Context) (string, time.Time, error) {
	sd := sm.getSessionDataFromContext(ctx)

	sd.mu.Lock()
	defer sd.mu.Unlock()

	if sd.sessionID == "" {
		var err error
		if sd.sessionID, err = createSessionID(); err != nil {
			return "", time.Time{}, err
		}
	}

	b, err := sm.Encode(sd.expiryTime, sd.values)
	if err != nil {
		return "", time.Time{}, err
	}

	expiry := sd.expiryTime

	if err := sm.Store.StoreCommit(ctx, sd.sessionID, b, expiry); err != nil {
		return "", time.Time{}, err
	}

	return sd.sessionID, expiry, nil
}

func (sm *SessionManager) WriteSessionCookie(ctx context.Context, w http.ResponseWriter,
	token string, expiry time.Time) {

	cookie := &http.Cookie{
		Name:     sm.CookieName,
		Value:    token,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Expires:  expiry,
	}

	if expiry.IsZero() {
		cookie.Expires = time.Unix(1, 0)
		cookie.MaxAge = -1
	}

	w.Header().Add("Set-Cookie", cookie.String())
	w.Header().Add("Cache-Control", `no-cache="Set-Cookie"`)
}

type sessionWriter struct {
	http.ResponseWriter
	request        *http.Request
	sessionManager *SessionManager
	written        bool
}

func (sw *sessionWriter) Write(b []byte) (int, error) {
	if !sw.written {
		sw.sessionManager.commitAndWriteSessionCookie(sw.ResponseWriter, sw.request)
		sw.written = true
	}

	return sw.ResponseWriter.Write(b)
}

func (sw *sessionWriter) WriteHeader(code int) {
	if !sw.written {
		sw.sessionManager.commitAndWriteSessionCookie(sw.ResponseWriter, sw.request)
		sw.written = true
	}

	sw.ResponseWriter.WriteHeader(code)
}
