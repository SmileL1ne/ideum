package post

import (
	"database/sql"
	"errors"
	"forum/internal/entity"
)

// TODO: Add new methods related to post manipulation (delete, update post)
type IPostRepository interface {
	SavePost(entity.PostCreateForm, int) (int, error)
	GetPost(int) (entity.PostEntity, error)
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

func (r *postRepository) SavePost(p entity.PostCreateForm, userID int) (int, error) {
	query := `
		INSERT INTO posts (title, content, user_id, created_at) 
		VALUES ($1, $2, $3, datetime('now', 'utc', '+12 hours'))`

	result, err := r.DB.Exec(query, p.Title, p.Content, userID)
	if err != nil {
		return 0, err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(postID), nil
}

func (r *postRepository) GetPost(postID int) (entity.PostEntity, error) {
	query := `
		SELECT p.id, p.title, p.content, p.created_at, u.username,
			SUM(CASE WHEN pr.is_like = true THEN 1 ELSE 0 END) as likes_count,
			SUM(CASE WHEN pr.is_like = false THEN 1 ELSE 0 END) as dislikes_count  
		FROM posts p
		INNER JOIN users u ON p.user_id = u.id
		LEFT JOIN post_reactions pr ON p.id = pr.post_id
		WHERE p.id=$1
		GROUP BY p.id
		`

	var post entity.PostEntity
	if err := r.DB.QueryRow(query, postID).Scan(&post.ID, &post.Title, &post.Content,
		&post.CreatedAt, &post.Username, &post.Likes, &post.Dislikes); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.PostEntity{}, entity.ErrNoRecord
		}
		return entity.PostEntity{}, err
	}

	return post, nil
}

func (r *postRepository) GetAllPosts() (*[]entity.PostEntity, error) {
	query := `
	SELECT p.id, p.title, p.content, p.created_at, u.username, 
		SUM(CASE WHEN pr.is_like = true THEN 1 ELSE 0 END) as likes_count,
		SUM(CASE WHEN pr.is_like = false THEN 1 ELSE 0 END) as dislikes_count
	FROM posts p
	INNER JOIN users u ON p.user_id = u.id
	LEFT JOIN post_reactions pr ON p.id = pr.post_id
	GROUP BY p.id
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}

	var posts []entity.PostEntity

	for rows.Next() {
		var post entity.PostEntity
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt,
			&post.Username, &post.Likes, &post.Dislikes); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &posts, nil
}
