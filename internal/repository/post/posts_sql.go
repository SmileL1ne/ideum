package repository

import (
	"database/sql"
	"forum/internal/entity"
)

type postRepo struct {
	*sql.DB
}

func New(db *sql.DB) *postRepo {
	return &postRepo{db}
}

// TODO: Implement this method
func (r *postRepo) SavePost(p entity.Post) error {
	return nil
}
