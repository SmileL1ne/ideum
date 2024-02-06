package service

import (
	"forum/internal/repository"
	"forum/internal/service/comments"
	"forum/internal/service/posts"
	"forum/internal/service/users"
)

type Services struct { // TODO: rename to 'Services'
	Post    posts.IPostService
	User    users.IUserService
	Comment comments.ICommentService
}

func New(r *repository.Repositories) *Services {
	return &Services{
		Post:    posts.NewPostsService(r.Posts),
		User:    users.NewUserService(r.Users),
		Comment: comments.NewCommentService(r.Comments),
	}
}
