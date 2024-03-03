package handlers

import (
	"fmt"
	"net/http"
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
func secureHeaders(next http.Handler) http.Handler {
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
