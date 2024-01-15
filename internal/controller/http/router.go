package http

import (
	"forum/internal/service"
	"html/template"
	"log/slog"
	"net/http"
	"os"
)

/*
	TODO:
	- Checkout comments in each registered handlers
	- Add static handler
	-- Add middlewares
*/

type routes struct {
	s         *service.Service
	l         *slog.Logger
	tempCache map[string]*template.Template
}

func NewRouter(l *slog.Logger, s *service.Service) http.Handler {
	router := http.NewServeMux()

	tempCache, err := newTemplateCache()
	if err != nil {
		l.Error(err.Error())
		os.Exit(1)
	}

	r := &routes{
		l:         l,
		s:         s,
		tempCache: tempCache,
	}

	router.HandleFunc("/static/", fileServer.ServeHTTP) // Add static handler

	router.HandleFunc("/", r.home)                   // Should be GET method
	router.HandleFunc("/posts/view", r.postView)     // Should be GET method
	router.HandleFunc("/posts/create", r.postCreate) // Should be GET method and redirect to postCreatePost if method is POST
	router.HandleFunc("/user/signup", r.userSignup)  // Should be GET method and redirect to userSignupPost if method is POST
	router.HandleFunc("/user/login", r.userLogin)    // Should be GET method and redirect to userLoginPost if method is POST
	router.HandleFunc("/user/logout", r.userLogout)  // Should be POST method

	return router
}
