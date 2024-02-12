package reaction

import (
	"database/sql"
	"errors"
	"forum/internal/entity"
)

type IReactionRepository interface {
	ExistsPost(int, int) (bool, error)
	ExistsComment(int, int) (bool, error)
	AddPost(bool, int, int) error
	AddComment(bool, int, int) error
	DeletePost(int, int) error
	DeleteComment(int, int) error
}

type reactionRepository struct {
	DB *sql.DB
}

func NewReactionRepo(db *sql.DB) *reactionRepository {
	return &reactionRepository{
		DB: db,
	}
}

var _ IReactionRepository = (*reactionRepository)(nil)

func (r *reactionRepository) ExistsPost(userID int, postID int) (bool, error) {
	var isLike bool

	query := `
		SELECT is_like
		FROM post_reactions
		WHERE user_id = $1 AND post_id = $2
	`

	err := r.DB.QueryRow(query, userID, postID).Scan(&isLike)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, entity.ErrNoRecord
		}
		return false, err
	}

	return isLike, nil
}

func (r *reactionRepository) AddPost(isLike bool, postID int, userID int) error {
	query := `
		INSERT INTO post_reactions (post_id, user_id, is_like, created_at)
		VALUES ($1, $2, $3, datetime('now', 'utc', '+12 hours'))
	`

	_, err := r.DB.Exec(query, postID, userID, isLike)

	return err
}

func (r *reactionRepository) DeletePost(postID int, userID int) error {
	query := `
		DELETE
		FROM post_reactions
		WHERE post_id = $1 AND user_id = $2
	`

	_, err := r.DB.Exec(query, postID, userID)

	return err
}

func (r *reactionRepository) ExistsComment(userID int, commentID int) (bool, error) {
	var isLike bool

	query := `
		SELECT is_like
		FROM comment_reactions
		WHERE user_id = $1 AND comment_id = $2
	`

	err := r.DB.QueryRow(query, userID, commentID).Scan(&isLike)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, entity.ErrNoRecord
		}
		return false, err
	}

	return isLike, nil
}

func (r *reactionRepository) AddComment(isLike bool, commentID int, userID int) error {
	query := `
		INSERT INTO comment_reactions (comment_id, user_id, is_like, created_at)
		VALUES ($1, $2, $3, datetime('now', 'utc', '+12 hours'))
	`

	_, err := r.DB.Exec(query, commentID, userID, isLike)

	return err
}

func (r *reactionRepository) DeleteComment(commentID int, userID int) error {
	query := `
		DELETE
		FROM comment_reactions
		WHERE comment_id = $1 AND user_id = $2
	`

	_, err := r.DB.Exec(query, commentID, userID)

	return err
}
