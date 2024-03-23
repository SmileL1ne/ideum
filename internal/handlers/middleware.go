package handlers

import (
	"fmt"
	"forum/internal/entity"
	"forum/pkg/rate"
	"net/http"
	"strings"
	"time"
)

// requireAuthentication middleware checks if user if authenticated
// and if not, redirects to login page.
//
// It also sets 'Cache-Control' to 'no-store' to avoid saving pages
// in browsers cache
func (r *Routes) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if !r.isAuthenticated(req) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		w.Header().Set("Cache-Control", "no-store")

		next.ServeHTTP(w, req)
	})
}

// secureHeaders middleware sets several headers to secure every response
func (r *Routes) secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self'")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, req)
	})
}

// recoverPanic middleware ensures that any panic WITHIN a handler would be
// handled by the anonymous function declared in this middleware
func (r *Routes) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				r.serverError(w, req, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, req)
	})
}

// limitRate middleware checks every request's rate, if it is more than limit that set it
// would block user for given amount of time (in seconds)
func (r *Routes) limitRate(next http.Handler) http.Handler {
	rl := rate.NewRateLimiter(time.Duration(r.cfg.RateInterval)*time.Second, r.cfg.RateLimit)

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if strings.HasPrefix(req.URL.Path, "/static/") {
			next.ServeHTTP(w, req)
			return
		}

		ip := req.RemoteAddr

		r.rateMu.Lock()
		user, ok := r.userRateLimits[ip]
		r.rateMu.Unlock()

		if ok {
			r.rateMu.Lock()

			if time.Since(user.lastReq) < user.penalty {
				r.rateMu.Unlock()
				r.rateLimitExceeded(w, user.penalty-time.Since(user.lastReq))
				return
			}
			r.rateMu.Unlock()
		}

		if rl.Limit() {
			r.updateRateLimit(ip, time.Now(), time.Duration(r.cfg.RatePenalty)*time.Second)
			r.rateLimitExceeded(w, time.Duration(r.cfg.RatePenalty)*time.Second)
			return
		}

		r.updateRateLimit(ip, time.Now(), 0)

		next.ServeHTTP(w, req)
	})
}

func (r *Routes) updateRateLimit(ip string, lastReq time.Time, penalty time.Duration) {
	r.rateMu.Lock()
	defer r.rateMu.Unlock()

	r.userRateLimits[ip] = userRateLimit{
		lastReq: lastReq,
		penalty: penalty,
	}
}

func (r *Routes) detectGuest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		role := r.sesm.GetUserRole(req.Context())
		if role == "" {
			r.sesm.PutUserRoleWithoutStatusChange(req.Context(), entity.GUEST)
		}

		next.ServeHTTP(w, req)
	})
}

func (r *Routes) requireAdminRights(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		userRole := r.sesm.GetUserRole(req.Context())

		if userRole != entity.ADMIN {
			r.forbidden(w)
			return
		}

		next.ServeHTTP(w, req)
	})
}
