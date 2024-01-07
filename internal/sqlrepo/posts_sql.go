package sqlrepo

import (
	"database/sql"
	"forum/internal/entity"
)

type PostsRepo struct {
	*sql.DB
}

func New(db *sql.DB) *PostsRepo {
	return &PostsRepo{db}
}

// TODO: Implement this method
func (r *PostsRepo) SavePost(p entity.Post) error {
	return nil
}
