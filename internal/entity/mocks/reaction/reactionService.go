package reaction

import (
	"errors"
	"forum/internal/entity"
	repo "forum/internal/repository/reaction"
	service "forum/internal/service/reaction"
)

type ReactionServiceMock struct {
	r repo.IReactionRepository
}

func NewReactionServiceMock(r repo.IReactionRepository) *ReactionServiceMock {
	return &ReactionServiceMock{
		r: r,
	}
}

var _ service.IReactionService = (*ReactionServiceMock)(nil)

func (rs *ReactionServiceMock) SetPostReaction(reaction string, postID int, userID int) error {
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
	isLikeDB, err := rs.r.ExistsPostReaction(userID, postID)
	if err != nil {
		if errors.Is(err, entity.ErrNoRecord) {
			return rs.r.AddPostReaction(isLike, postID, userID)
		}
		return err
	}

	err = rs.r.DeletePostReaction(postID, userID)
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
		return rs.r.AddPostReaction(isLike, postID, userID)
	}
}

// Same principle to reactions handling in posts
func (rs *ReactionServiceMock) SetCommentReaction(reaction string, commentID int, userID int) error {
	var isLike bool
	switch reaction {
	case "like":
		isLike = true
	case "dislike":
		isLike = false
	default:
		return entity.ErrInvalidURLPath
	}

	isLikeDB, err := rs.r.ExistsCommentReaction(userID, commentID)
	if err != nil {
		if errors.Is(err, entity.ErrNoRecord) {
			return rs.r.AddCommentReaction(isLike, commentID, userID)
		}
		return err
	}

	err = rs.r.DeleteCommentReaction(commentID, userID)
	if err != nil {
		return err
	}

	if isLike == isLikeDB {
		return err
	} else {
		return rs.r.AddCommentReaction(isLike, commentID, userID)
	}
}
