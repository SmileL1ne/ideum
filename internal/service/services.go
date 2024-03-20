package service

import (
	"forum/internal/repository"
	"forum/internal/service/comment"
	"forum/internal/service/image"
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
	Image    image.IImageService
}

func New(r *repository.Repositories) *Services {
	return &Services{
		Post:     post.NewPostsService(r.Post, image.NewImageService(r.Image), tag.NewTagService(r.Tag), comment.NewCommentService(r.Comment, user.NewUserService(r.User)), user.NewUserService(r.User)),
		User:     user.NewUserService(r.User),
		Comment:  comment.NewCommentService(r.Comment, user.NewUserService(r.User)),
		Reaction: reaction.NewReactionService(r.Reaction),
		Tag:      tag.NewTagService(r.Tag),
		Image:    image.NewImageService(r.Image),
	}
}
