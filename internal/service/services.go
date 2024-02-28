package service

import (
	"forum/internal/repository"
	"forum/internal/service/comment"
	"forum/internal/service/post"
	"forum/internal/service/reaction"
	"forum/internal/service/tag"
	"forum/internal/service/user"
)

type Services struct {
	Post     post.IPostService
	User     user.IUserService
	Comment  comment.ICommentService
	Reaction reaction.IReactionService
	Tag      tag.ITagService
}

func New(r *repository.Repositories) *Services {
	return &Services{
		Post:     post.NewPostsService(r.Post),
		User:     user.NewUserService(r.User),
		Comment:  comment.NewCommentService(r.Comment),
		Reaction: reaction.NewReactionService(r.Reaction),
		Tag:      tag.NewTagService(r.Tag),
	}
}
