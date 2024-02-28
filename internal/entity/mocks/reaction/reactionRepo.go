package reaction

import (
	"forum/internal/entity"
	"forum/internal/repository/reaction"
)

type ReactionRepoMock struct{}

func NewReactionRepoMock() *ReactionRepoMock {
	return &ReactionRepoMock{}
}

var _ reaction.IReactionRepository = (*ReactionRepoMock)(nil)

func (r *ReactionRepoMock) ExistsPostReaction(userID int, postID int) (bool, error) {
	if userID == 1 && postID == 1 {
		return true, nil
	} else {
		return false, entity.ErrNoRecord
	}
}

func (r *ReactionRepoMock) AddPostReaction(isLike bool, postID int, userID int) error {
	return nil
}

func (r *ReactionRepoMock) DeletePostReaction(postID int, userID int) error {
	return nil
}

func (r *ReactionRepoMock) ExistsCommentReaction(userID int, commentID int) (bool, error) {
	if userID == 1 && commentID == 1 {
		return true, nil
	} else {
		return false, entity.ErrNoRecord
	}
}

func (r *ReactionRepoMock) AddCommentReaction(isLike bool, commentID int, userID int) error {
	return nil
}

func (r *ReactionRepoMock) DeleteCommentReaction(commentID int, userID int) error {
	return nil
}
