package http

import (
	"forum/internal/service"
	"html/template"
	"log/slog"
	"net/http"
	"os"
)

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

	router.HandleFunc("/", r.home)

	return router
}
