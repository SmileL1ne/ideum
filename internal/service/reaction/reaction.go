package reaction

import (
	"errors"
	"forum/internal/entity"
	"forum/internal/repository/reaction"
	"strconv"
)

type IReactionService interface {
	AddOrDeletePost(bool, string, int) error
	AddOrDeleteComment(bool, string, int) error
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

func (rs *reactionService) AddOrDeletePost(isLike bool, postIDStr string, userID int) error {
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		return entity.ErrInvalidPostId
	}

	// Check if reaction by user exists in table, if not it would return
	// entity.ErrNoRecord error, otherwise it would return reaction left by
	// that user
	isLikeDB, err := rs.reactsRepo.ExistsPost(userID)
	if err != nil {
		if errors.Is(err, entity.ErrNoRecord) {
			return rs.reactsRepo.AddPost(isLike, postID, userID)
		}
		return err
	}

	err = rs.reactsRepo.DeletePost(postID, userID)
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
		return rs.reactsRepo.AddPost(isLike, postID, userID)
	}
}

// Same principle to reactions handling in posts
func (rs *reactionService) AddOrDeleteComment(isLike bool, commentIDStr string, userID int) error {
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		return entity.ErrInvalidCommentId
	}

	isLikeDB, err := rs.reactsRepo.ExistsComment(userID)
	if err != nil {
		if errors.Is(err, entity.ErrNoRecord) {
			return rs.reactsRepo.AddComment(isLike, commentID, userID)
		}
		return err
	}

	err = rs.reactsRepo.DeleteComment(commentID, userID)
	if err != nil {
		return err
	}

	if isLike == isLikeDB {
		return err
	} else {
		return rs.reactsRepo.AddComment(isLike, commentID, userID)
	}
}
