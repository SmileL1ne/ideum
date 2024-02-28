package reaction

import (
	"database/sql"
	"errors"
	"forum/internal/entity"
)

type IReactionRepository interface {
	ExistsPostReaction(int, int) (bool, error)
	ExistsCommentReaction(int, int) (bool, error)
	AddPostReaction(bool, int, int) error
	AddCommentReaction(bool, int, int) error
	DeletePostReaction(int, int) error
	DeleteCommentReaction(int, int) error
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

func (r *reactionRepository) ExistsPostReaction(userID int, postID int) (bool, error) {
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

func (r *reactionRepository) AddPostReaction(isLike bool, postID int, userID int) error {
	query := `
		INSERT INTO post_reactions (post_id, user_id, is_like, created_at)
		VALUES ($1, $2, $3, datetime('now', 'localtime'))
	`

	_, err := r.DB.Exec(query, postID, userID, isLike)

	return err
}

func (r *reactionRepository) DeletePostReaction(postID int, userID int) error {
	query := `
		DELETE
		FROM post_reactions
		WHERE post_id = $1 AND user_id = $2
	`

	_, err := r.DB.Exec(query, postID, userID)

	return err
}

func (r *reactionRepository) ExistsCommentReaction(userID int, commentID int) (bool, error) {
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

func (r *reactionRepository) AddCommentReaction(isLike bool, commentID int, userID int) error {
	query := `
		INSERT INTO comment_reactions (comment_id, user_id, is_like, created_at)
		VALUES ($1, $2, $3, datetime('now', 'localtime'))
	`

	_, err := r.DB.Exec(query, commentID, userID, isLike)

	return err
}

func (r *reactionRepository) DeleteCommentReaction(commentID int, userID int) error {
	query := `
		DELETE
		FROM comment_reactions
		WHERE comment_id = $1 AND user_id = $2
	`

	_, err := r.DB.Exec(query, commentID, userID)

	return err
}
