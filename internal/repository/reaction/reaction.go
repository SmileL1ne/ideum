package reaction

import (
	"database/sql"
	"fmt"
)

type IReactionRepository interface {
	Save(string, int, int) error
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

func (r *reactionRepository) Save(reaction string, postID int, userID int) error {
	query := fmt.Sprintf(`
		UPDATE post_reactions
		SET %s = %s + 1, 
			updated_at = datetime('now', 'utc', '+12 hours')
		WHERE post_id = $1 AND user_id = $2
	`, reaction, reaction)

	_, err := r.DB.Exec(query, postID, userID)
	if err != nil {
		return err
	}

	return err
}
