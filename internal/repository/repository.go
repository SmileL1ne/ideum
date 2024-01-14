package repository

import (
	"database/sql"
	"forum/internal/repository/posts"
)

type Repository struct {
	Posts posts.PostRepository
}

func New(db *sql.DB) *Repository {
	return &Repository{
		Posts: posts.NewPostRepo(db),
	}
}
