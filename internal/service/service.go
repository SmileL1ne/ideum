package service

import (
	"forum/internal/repository"
	"forum/internal/service/posts"
	"forum/internal/service/users"
)

type Service struct {
	Post posts.PostService
	User users.UserService
}

func New(r *repository.Repository) *Service {
	return &Service{
		Post: posts.NewPostsService(r.Posts),
		User: users.NewUserService(r.Users),
	}
}
