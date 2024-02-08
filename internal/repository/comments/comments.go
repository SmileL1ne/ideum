package comments

import (
	"database/sql"
	"forum/internal/entity"
)

type ICommentRepository interface {
	SaveComment(entity.CommentCreateForm, int) error
	GetAllCommentsForPost(int) (*[]entity.CommentEntity, error)
}

type commentRepository struct {
	DB *sql.DB
}

var _ ICommentRepository = (*commentRepository)(nil)

func NewCommentRepository(db *sql.DB) *commentRepository {
	return &commentRepository{
		DB: db,
	}
}

func (r *commentRepository) SaveComment(c entity.CommentCreateForm, postId int) error {
	query := `INSERT INTO comments (content, post_id, created_at) VALUES ($1, $2, datetime('now', 'utc', '+12 hours'))`

	_, err := r.DB.Exec(query, c.Content, postId)

	return err
}

func (r *commentRepository) GetAllCommentsForPost(postId int) (*[]entity.CommentEntity, error) {
	query := `
		SELECT * 
		FROM comments
		WHERE post_id = $1
		`

	rows, err := r.DB.Query(query, postId)
	if err != nil {
		return nil, err
	}

	var comments []entity.CommentEntity

	for rows.Next() {
		var comment entity.CommentEntity
		if err := rows.Scan(&comment.Id, &comment.Content, &comment.CreatedAt, &comment.PostID); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &comments, nil
}
