package reaction

import (
	"errors"
	"forum/internal/entity"
	"forum/internal/repository/reaction"
)

type IReactionService interface {
	SetPostReaction(string, int, int) error
	SetCommentReaction(string, int, int) error
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

func (rs *reactionService) SetPostReaction(reaction string, postID int, userID int) error {
	var isLike bool
	switch reaction {
	case "like":
		isLike = true
	case "dislike":
		isLike = false
	default:
		return entity.ErrInvalidURLPath
	}

	// Check if reaction by user exists in table, if not it would return
	// entity.ErrNoRecord error, otherwise it would return reaction left by
	// that user
	isLikeDB, err := rs.reactsRepo.ExistsPostReaction(userID, postID)
	if err != nil {
		if errors.Is(err, entity.ErrNoRecord) {
			return rs.reactsRepo.AddPostReaction(isLike, postID, userID)
		}
		return err
	}

	err = rs.reactsRepo.DeletePostReaction(postID, userID)
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
		return rs.reactsRepo.AddPostReaction(isLike, postID, userID)
	}
}

// Same principle to reactions handling in posts
func (rs *reactionService) SetCommentReaction(reaction string, commentID int, userID int) error {
	var isLike bool
	switch reaction {
	case "like":
		isLike = true
	case "dislike":
		isLike = false
	default:
		return entity.ErrInvalidURLPath
	}

	isLikeDB, err := rs.reactsRepo.ExistsCommentReaction(userID, commentID)
	if err != nil {
		if errors.Is(err, entity.ErrNoRecord) {
			return rs.reactsRepo.AddCommentReaction(isLike, commentID, userID)
		}
		return err
	}

	err = rs.reactsRepo.DeleteCommentReaction(commentID, userID)
	if err != nil {
		return err
	}

	if isLike == isLikeDB {
		return err
	} else {
		return rs.reactsRepo.AddCommentReaction(isLike, commentID, userID)
	}
}
