package comment

import (
	"database/sql"
	"forum/internal/entity"
)

type ICommentRepository interface {
	SaveComment(entity.CommentCreateForm, int, int) error
	GetAllCommentsForPost(int) (*[]entity.CommentEntity, error)
	ExistsComment(int) (bool, error)
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

func (r *commentRepository) SaveComment(c entity.CommentCreateForm, postID int, userID int) error {
	query := `
		INSERT INTO comments (content, post_id, user_id, created_at) 
		VALUES ($1, $2, $3, datetime('now', 'localtime'))
		`

	_, err := r.DB.Exec(query, c.Content, postID, userID)

	return err
}

func (r *commentRepository) GetAllCommentsForPost(postID int) (*[]entity.CommentEntity, error) {
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

func (r *commentRepository) ExistsComment(commentID int) (bool, error) {
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
