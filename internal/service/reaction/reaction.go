package reaction

import (
	"errors"
	"forum/internal/entity"
	"forum/internal/repository/reaction"
	"strconv"
)

type IReactionService interface {
	AddOrDelete(bool, string, int) error
}

type reactionService struct {
	reactsRepo reaction.IReactionRepository
}

func NewReactionService(r reaction.IReactionRepository) *reactionService {
	return &reactionService{
		reactsRepo: r,
	}
}

var _ IReactionService = (*reactionService)(nil)

func (rs *reactionService) AddOrDelete(isLike bool, postIDStr string, userID int) error {
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		return entity.ErrInvalidPostId
	}

	// Check if reaction by user exists in table, if not it would return
	// entity.ErrNoRecord error, otherwise it would return reaction left by
	// that user
	isLikeDB, err := rs.reactsRepo.Exists(userID)
	if err != nil {
		if errors.Is(err, entity.ErrNoRecord) {
			return rs.reactsRepo.Add(isLike, postID, userID)
		}
		return err
	}

	err = rs.reactsRepo.Delete(postID, userID)
	if err != nil {
		return err
	}

	// If new reaction is the same as was in table just return an error
	// of deleted row.
	//
	// If not, add new reaction to the table (overall, replacing old reaction with new one)
	if isLike == isLikeDB {
		return err
	} else {
		return rs.reactsRepo.Add(isLike, postID, userID)
	}
}
