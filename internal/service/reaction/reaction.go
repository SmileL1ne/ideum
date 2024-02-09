package reaction

import (
	"forum/internal/entity"
	"forum/internal/repository/reaction"
	"strconv"
)

type IReactionService interface {
	Save(bool, string, int) error
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

func (rs *reactionService) Save(isLike bool, postIDStr string, userID int) error {
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		return entity.ErrInvalidPostId
	}

	reaction := ""
	if isLike {
		reaction = "likes"
	} else {
		reaction = "dislikes"
	}

	return rs.reactsRepo.Save(reaction, postID, userID)
}
