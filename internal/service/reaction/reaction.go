package reaction

import (
	"errors"
	"forum/internal/entity"
	"forum/internal/repository/reaction"
)

type IReactionService interface {
	AddOrDeletePost(string, int, int) error
	AddOrDeleteComment(string, int, int) error
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

func (rs *reactionService) AddOrDeletePost(reaction string, postID int, userID int) error {
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
	isLikeDB, err := rs.reactsRepo.ExistsPost(userID, postID)
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
func (rs *reactionService) AddOrDeleteComment(reaction string, commentID int, userID int) error {
	var isLike bool
	switch reaction {
	case "like":
		isLike = true
	case "dislike":
		isLike = false
	default:
		return entity.ErrInvalidURLPath
	}

	isLikeDB, err := rs.reactsRepo.ExistsComment(userID, commentID)
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
