package http

import (
	"forum/internal/service"
	"html/template"
	"log/slog"
	"net/http"
	"os"
)

type routes struct {
	p         service.PostService
	l         *slog.Logger
	tempCache map[string]*template.Template
}

func NewRouter(l *slog.Logger, p service.PostService) http.Handler {
	router := http.NewServeMux()

	tempCache, err := newTemplateCache()
	if err != nil {
		l.Error(err.Error())
		os.Exit(1)
	}

	r := &routes{
		l:         l,
		p:         p,
		tempCache: tempCache,
	}

	router.HandleFunc("/", r.home)

	return router
}
