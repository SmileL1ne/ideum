package posts

import (
	"database/sql"
	"errors"
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
	stmt := `SELECT * FROM posts WHERE id=?`

	post := entity.Post{}
	if err := r.DB.QueryRow(stmt, postId).Scan(&post.Id, &post.Title, &post.Content, &post.Created); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Post{}, entity.ErrNoRecord
		}
		return entity.Post{}, err
	}

	return post, nil
}

func (r *postRepository) GetAllPosts() ([]entity.Post, error) {
	return nil, nil
}
