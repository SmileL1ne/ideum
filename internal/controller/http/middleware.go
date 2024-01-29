package http

import "net/http"

func (r *routes) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if !r.isAuthenticated(req) {
			http.Redirect(w, req, "/user/login", http.StatusSeeOther)
			return
		}

		w.Header().Set("Cache-Control", "no-store")

		next.ServeHTTP(w, req)
	})
}

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self'")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, req)
	})
}
