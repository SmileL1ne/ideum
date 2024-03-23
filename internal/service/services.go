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
	postService := post.NewPostsService(r.Post, image.NewImageService(r.Image), tag.NewTagService(r.Tag), comment.NewCommentService(r.Comment, user.NewUserService(r.User)), user.NewUserService(r.User))
	commentService := comment.NewCommentService(r.Comment, user.NewUserService(r.User))
	userService := user.NewUserService(r.User)
	return &Services{
		Post:     postService,
		User:     userService,
		Comment:  commentService,
		Reaction: reaction.NewReactionService(r.Reaction, postService, commentService, userService),
		Tag:      tag.NewTagService(r.Tag),
		Image:    image.NewImageService(r.Image),
	}
}
