package repository

import (
	"database/sql"
	"forum/internal/repository/comment"
	"forum/internal/repository/post"
	"forum/internal/repository/reaction"
	"forum/internal/repository/user"
)

type Repositories struct {
	Post     post.IPostRepository
	User     user.IUserRepository
	Comment  comment.ICommentRepository
	Reaction reaction.IReactionRepository
}

func New(db *sql.DB) *Repositories {
	return &Repositories{
		Post:     post.NewPostRepo(db),
		User:     user.NewUserRepo(db),
		Comment:  comment.NewCommentRepo(db),
		Reaction: reaction.NewReactionRepo(db),
	}
}
