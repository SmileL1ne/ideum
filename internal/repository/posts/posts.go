package posts

import (
	"database/sql"
	"forum/internal/entity"
)

// TODO: Add new methods related to post manipulation (delete, update post)
type PostRepository interface {
	SavePost(entity.Post) (int, error)
	GetPost(postId int) (entity.Post, error)
	GetAllPosts() ([]entity.Post, error)
}

type postRepository struct {
	DB *sql.DB
}

func NewPostRepo(db *sql.DB) *postRepository {
	return &postRepository{
		DB: db,
	}
}

// TODO: Implement this method
func (r *postRepository) SavePost(p entity.Post) (int, error) {
	stmt := `INSERT INTO posts (title, content, created) VALUES (?, ?, datetime('now', 'utc', '+12 hours'))`

	result, err := r.DB.Exec(stmt, p.Title, p.Content)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}

	return int(id), nil
}

func (r *postRepository) GetPost(postId int) (entity.Post, error) {
	return entity.Post{}, nil
}

func (r *postRepository) GetAllPosts() ([]entity.Post, error) {
	return nil, nil
}
