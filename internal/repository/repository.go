package repository

import (
	"database/sql"
	"forum/internal/repository/posts"
	"forum/internal/repository/users"
)

type Repository struct {
	Posts posts.PostRepository
	Users users.UserRepository
}

func New(db *sql.DB) *Repository {
	return &Repository{
		Posts: posts.NewPostRepo(db),
		Users: users.NewUserRepo(db),
	}
}
