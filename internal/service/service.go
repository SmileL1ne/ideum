package service

import (
	"forum/internal/repository"
	"forum/internal/service/posts"
)

type Service struct {
	Post posts.PostService
}

func New(r *repository.Repository) *Service {
	return &Service{
		Post: posts.NewPostsService(r.Posts),
	}
}
