package comment

import (
	"database/sql"
	"errors"
	"forum/internal/entity"
)

type ICommentRepository interface {
	Insert(entity.CommentCreateForm, int, int) error
	GetAllForPost(int) (*[]entity.CommentEntity, error)
	GetAllUserCommentsForPost(userID, postID int) (*[]entity.CommentEntity, error)
	Exists(int) (bool, error)
	Delete(commentID, userID int) error
	DeleteByPrivileged(commentID int) error
	GetAuthorID(commentID int) (int, error)
}

type commentRepository struct {
	DB *sql.DB
}

var _ ICommentRepository = (*commentRepository)(nil)

func NewCommentRepo(db *sql.DB) *commentRepository {
	return &commentRepository{
		DB: db,
	}
}

func (r *commentRepository) Insert(c entity.CommentCreateForm, postID int, userID int) error {
	query := `
		INSERT INTO comments (content, post_id, user_id, created_at) 
		VALUES ($1, $2, $3, datetime('now', 'localtime'))
		`

	_, err := r.DB.Exec(query, c.Content, postID, userID)

	return err
}

func (r *commentRepository) GetAllForPost(postID int) (*[]entity.CommentEntity, error) {
	query := `
		SELECT c.id, c.content, c.created_at, c.post_id, u.username, 
			SUM(CASE WHEN cr.is_like = true THEN 1 ELSE 0 END) as likes_count,
			SUM(CASE WHEN cr.is_like = false THEN 1 ELSE 0 END) as dislikes_count
		FROM comments c
		INNER JOIN users u ON c.user_id = u.id
		LEFT JOIN comment_reactions cr ON c.id = cr.comment_id
		WHERE c.post_id = $1
		GROUP BY c.id
	`

	rows, err := r.DB.Query(query, postID)
	if err != nil {
		return nil, err
	}

	var comments []entity.CommentEntity

	for rows.Next() {
		var comment entity.CommentEntity
		if err := rows.Scan(&comment.ID, &comment.Content, &comment.CreatedAt,
			&comment.PostID, &comment.Username, &comment.Likes, &comment.Dislikes); err != nil {

			return nil, err
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &comments, nil
}

func (r *commentRepository) GetAllUserCommentsForPost(userID, postID int) (*[]entity.CommentEntity, error) {
	query := `
		SELECT c.id, c.content, c.created_at, c.post_id, u.username,
			SUM(CASE WHEN cr.is_like = true THEN 1 ELSE 0 END) as likes_count,
			SUM(CASE WHEN cr.is_like = false THEN 1 ELSE 0 END) as dislikes_count
		FROM comments c
		INNER JOIN users u ON c.user_id = u.id
		LEFT JOIN comment_reactions cr ON c.id = cr.comment_id
		WHERE c.user_id = $1 AND c.post_id = $2
		GROUP BY c.id
	`

	rows, err := r.DB.Query(query, userID, postID)
	if err != nil {
		return nil, err
	}

	var comments []entity.CommentEntity

	for rows.Next() {
		var comment entity.CommentEntity
		if err := rows.Scan(&comment.ID, &comment.Content, &comment.CreatedAt,
			&comment.PostID, &comment.Username, &comment.Likes, &comment.Dislikes); err != nil {

			return nil, err
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &comments, nil
}

func (r *commentRepository) Exists(commentID int) (bool, error) {
	var exists bool

	query := `
		SELECT EXISTS(
			SELECT true
			FROM comments
			WHERE id = $1
		)
	`

	err := r.DB.QueryRow(query, commentID).Scan(&exists)
	return exists, err
}

func (r *commentRepository) Delete(commentID, userID int) error {
	query := `
		DELETE FROM comments
		WHERE id = $1 AND user_id = $2
	`

	_, err := r.DB.Exec(query, commentID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.ErrForbiddenAccess
		}
		return err
	}

	return nil
}

func (r *commentRepository) DeleteByPrivileged(commentID int) error {
	query := `
		DELETE FROM comments
		WHERE id = $1
	`

	_, err := r.DB.Exec(query, commentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.ErrCommentNotFound
		}
		return err
	}

	return nil
}

func (r *commentRepository) GetAuthorID(commentID int) (int, error) {
	query := `
		SELECT c.user_id
		FROM comments c
		WHERE c.id = $1
	`

	var userID int

	err := r.DB.QueryRow(query, commentID).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, entity.ErrCommentNotFound
		}
		return 0, err
	}

	return userID, nil
}
