package service

import (
	"forum/internal/repository"
	"forum/internal/service/comment"
	"forum/internal/service/post"
	"forum/internal/service/user"
)

type Services struct { // TODO: rename to 'Services'
	Post    post.IPostService
	User    user.IUserService
	Comment comment.ICommentService
}

func New(r *repository.Repositories) *Services {
	return &Services{
		Post:    post.NewPostsService(r.Posts),
		User:    user.NewUserService(r.Users),
		Comment: comment.NewCommentService(r.Comments),
	}
}
