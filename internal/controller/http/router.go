package http

import (
	"forum/internal/usecase"
	"log/slog"
	"net/http"
)

func NewRouter(l *slog.Logger, p *usecase.PostsUseCase) http.Handler {
	router := http.NewServeMux()
	pr := routes{
		l: l,
		p: p,
	}

	router.HandleFunc("/", pr.home)

	return router
}
