package http

import (
	"forum/internal/entity"
	"forum/internal/usecase"
	"log/slog"
	"net/http"
)

/*
	TODO:
	- Finish home handler
*/

type routes struct {
	p *usecase.PostsUseCase
	l *slog.Logger
}

func (r routes) home(w http.ResponseWriter, req *http.Request) {
	info := entity.Post{
		Title: "naaaaaaaaaaaaaaaaaaah",
	}
	err := r.p.MakeNewPost(info)
	if err != nil {
		panic(err)
	}
	w.Write([]byte(info.Title))
}
