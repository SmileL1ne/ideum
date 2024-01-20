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
	query := `INSERT INTO posts (title, content, created) VALUES ($1, $2, datetime('now', 'utc', '+12 hours'))`

	result, err := r.DB.Exec(query, p.Title, p.Content)
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
	query := `SELECT * FROM posts WHERE id=$1`

	var post entity.Post
	if err := r.DB.QueryRow(query, postId).Scan(&post.Id, &post.Title, &post.Content, &post.Created); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Post{}, entity.ErrNoRecord
		}
		return entity.Post{}, err
	}

	return post, nil
}

func (r *postRepository) GetAllPosts() ([]entity.Post, error) {
	query := `SELECT * FROM posts`

	var posts []entity.Post
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var post entity.Post
		if err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.Created); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
