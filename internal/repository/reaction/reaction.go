package reaction

import (
	"database/sql"
	"errors"
	"forum/internal/entity"
)

type IReactionRepository interface {
	Exists(int) (bool, error)
	Add(bool, int, int) error
	Delete(int, int) error
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

func (r *reactionRepository) Exists(userID int) (bool, error) {
	var isLike bool

	// query := `
	// 	SELECT EXISTS (
	// 		SELECT true
	// 		FROM post_reactions
	// 		WHERE user_id = $1
	// 	)
	// `

	query := `
		SELECT is_like
		FROM post_reactions
		WHERE user_id = $1
	`

	err := r.DB.QueryRow(query, userID).Scan(&isLike)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, entity.ErrNoRecord
		}
		return false, err
	}

	return isLike, nil
}

func (r *reactionRepository) Add(isLike bool, postID int, userID int) error {
	query := `
		INSERT INTO post_reactions (post_id, user_id, is_like, created_at)
		VALUES ($1, $2, $3, datetime('now', 'utc', '+12 hours'))
	`

	_, err := r.DB.Exec(query, postID, userID, isLike)

	return err
}

func (r *reactionRepository) Delete(postID int, userID int) error {
	query := `
		DELETE
		FROM post_reactions
		WHERE post_id = $1 AND user_id = $2
	`

	_, err := r.DB.Exec(query, postID, userID)

	return err
}
