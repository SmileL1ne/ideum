package repository

import (
	"database/sql"
	"forum/internal/repository/comments"
	"forum/internal/repository/posts"
	"forum/internal/repository/users"
)

type Repositories struct {
	Posts    posts.IPostRepository
	Users    users.IUserRepository
	Comments comments.ICommentRepository
}

func New(db *sql.DB) *Repositories {
	return &Repositories{
		Posts:    posts.NewPostRepo(db),
		Users:    users.NewUserRepo(db),
		Comments: comments.NewCommentRepository(db),
	}
}
