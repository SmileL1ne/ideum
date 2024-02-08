package repository

import (
	"database/sql"
	"forum/internal/repository/comment"
	"forum/internal/repository/post"
	"forum/internal/repository/user"
)

type Repositories struct {
	Posts    post.IPostRepository
	Users    user.IUserRepository
	Comments comment.ICommentRepository
}

func New(db *sql.DB) *Repositories {
	return &Repositories{
		Posts:    post.NewPostRepo(db),
		Users:    user.NewUserRepo(db),
		Comments: comment.NewCommentRepository(db),
	}
}
