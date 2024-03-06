package repository

import (
	"database/sql"
	"forum/internal/repository/comment"
	"forum/internal/repository/image"
	"forum/internal/repository/post"
	"forum/internal/repository/reaction"
	"forum/internal/repository/tag"
	"forum/internal/repository/user"
)

type Repositories struct {
	Post     post.IPostRepository
	User     user.IUserRepository
	Comment  comment.ICommentRepository
	Reaction reaction.IReactionRepository
	Tag      tag.ITagRepository
	Image    image.IImageRepository
}

func New(db *sql.DB) *Repositories {
	return &Repositories{
		Post:     post.NewPostRepo(db),
		User:     user.NewUserRepo(db),
		Comment:  comment.NewCommentRepo(db),
		Reaction: reaction.NewReactionRepo(db),
		Tag:      tag.NewTagRepo(db),
		Image:    image.NewImageRepo(db),
	}
}
