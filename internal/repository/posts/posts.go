package posts

import (
	"database/sql"
	"forum/internal/entity"
)

// TODO: Add new methods related to post manipulation (delete, update post)
type PostRepository interface {
	GetPost(postId int) (entity.Post, error)
	GetAllPosts() ([]entity.Post, error)
	SavePost(entity.Post) error
}

type postRepository struct {
	DB *sql.DB
}

func NewPostRepo(db *sql.DB) *postRepository {
	return &postRepository{
		DB: db,
	}
}

func (r *postRepository) GetPost(postId int) (entity.Post, error) {
	return entity.Post{}, nil
}

func (r *postRepository) GetAllPosts() ([]entity.Post, error) {
	return nil, nil
}

// TODO: Implement this method
func (r *postRepository) SavePost(p entity.Post) error {
	return nil
}
