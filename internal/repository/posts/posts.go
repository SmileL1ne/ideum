package posts

import (
	"database/sql"
	"errors"
	"forum/internal/entity"
)

// TODO: Add new methods related to post manipulation (delete, update post)
type IPostRepository interface {
	SavePost(entity.PostCreateForm) (int, error)
	GetPost(postId int) (entity.PostEntity, error)
	GetAllPosts() (*[]entity.PostEntity, error)
}

type postRepository struct {
	DB *sql.DB
}

func NewPostRepo(db *sql.DB) *postRepository {
	return &postRepository{
		DB: db,
	}
}

var _ IPostRepository = (*postRepository)(nil)

func (r *postRepository) SavePost(p entity.PostCreateForm) (int, error) {
	query := `INSERT INTO posts (title, content, created) VALUES ($1, $2, datetime('now', 'utc', '+12 hours'))`

	result, err := r.DB.Exec(query, p.Title, p.Content)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *postRepository) GetPost(postId int) (entity.PostEntity, error) {
	query := `SELECT * FROM posts WHERE id=$1`

	var post entity.PostEntity
	if err := r.DB.QueryRow(query, postId).Scan(&post.Id, &post.Title, &post.Content, &post.Created); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.PostEntity{}, entity.ErrNoRecord
		}
		return entity.PostEntity{}, err
	}

	return post, nil
}

func (r *postRepository) GetAllPosts() (*[]entity.PostEntity, error) {
	query := `SELECT * FROM posts`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}

	var posts []entity.PostEntity

	for rows.Next() {
		var post entity.PostEntity
		if err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.Created); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &posts, nil
}
