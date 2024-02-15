package sesm

import (
	"context"
	"forum/pkg/sesm/sqlit3store"
	"log"
	"net/http"
	"time"
)

/*
	Sesm is a session manager that helps to manage user sessions
	by providing convinient tools for it.
*/

type SessionManager struct {
	Store      *sqlit3store.SQLite3Store
	Lifetime   time.Duration
	CookieName string
	ContextKey contextKey
}

// New returns pointer to new SessionManager struct
func New() *SessionManager {
	return &SessionManager{
		Lifetime:   12 * time.Hour,
		ContextKey: generateContextKey(),
		CookieName: "session",
	}
}

// LoadAndSave is a middleware that loads session from cookie puts it
// to the request's context.
//
// If no session exists, it creates new automatically creates a new one.
func (sm *SessionManager) LoadAndSave(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Vary", "Cookie")

		var sessionID string
		cookie, err := req.Cookie(sm.CookieName)
		if err == nil {
			sessionID = cookie.Value
		}

		ctx, err := sm.Load(req.Context(), sessionID)
		if err != nil {
			log.Println("LoadAndSave:", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		sessionReq := req.WithContext(ctx)

		sessionWriter := &sessionWriter{
			ResponseWriter: w,
			request:        sessionReq,
			sessionManager: sm,
		}

		next.ServeHTTP(sessionWriter, sessionReq)

		// Commit changed data and write it to cookie if not by the
		// end of handler usage
		if !sessionWriter.written {
			sm.commitAndWriteSessionCookie(sessionWriter, sessionReq)
		}
	})
}

// commitAndWriteSessionCookie commits changes to database and saves it into cookie.
//
// It does it only in case of data being modified.
func (sm *SessionManager) commitAndWriteSessionCookie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch sm.Status(ctx) {
	case Modified:
		token, expiry, err := sm.commit(ctx)
		if err != nil {
			log.Print("Commit session:", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		sm.writeSessionCookie(ctx, w, token, expiry)
	case Destroyed:
		sm.writeSessionCookie(ctx, w, "", time.Time{})
	}

}

// commit attempts to save changes in database.
//
// If given context doesn't hold session it creates new one.
//
// It retures saved token, expiration date and error (if encountered).
func (sm *SessionManager) commit(ctx context.Context) (string, time.Time, error) {
	sd := sm.getSessionDataFromContext(ctx)

	sd.mu.Lock()
	defer sd.mu.Unlock()

	if sd.sessionID == "" {
		var err error
		if sd.sessionID, err = createSessionID(); err != nil {
			return "", time.Time{}, err
		}
	}

	b, err := sm.encode(sd.expiryTime, sd.values)
	if err != nil {
		return "", time.Time{}, err
	}

	expiry := sd.expiryTime

	if err := sm.Store.StoreCommit(ctx, sd.sessionID, b, expiry); err != nil {
		return "", time.Time{}, err
	}

	return sd.sessionID, expiry, nil
}

// writeSessionCookie creates new cookie and and saves it in response.
func (sm *SessionManager) writeSessionCookie(ctx context.Context, w http.ResponseWriter,
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

// sessionWriter overrides ResponseWriter's methods to save changes
// in session data before executing Write or WriteHeader methods
// (essentially, when handler finishes it's work)
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
