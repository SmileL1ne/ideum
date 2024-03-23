package reaction

import (
	"errors"
	"forum/internal/entity"
	"forum/internal/repository/reaction"
	"forum/internal/service/comment"
	"forum/internal/service/post"
	"forum/internal/service/user"
)

type IReactionService interface {
	SetPostReaction(reaction string, postID, userID int) error
	SetCommentReaction(reaction string, commentID, postID, userID int) error
}

type reactionService struct {
	reactsRepo     reaction.IReactionRepository
	postService    post.IPostService
	commentService comment.ICommentService
	userService    user.IUserService
}

func NewReactionService(r reaction.IReactionRepository, p post.IPostService, c comment.ICommentService, u user.IUserService) *reactionService {
	return &reactionService{
		reactsRepo:     r,
		postService:    p,
		commentService: c,
		userService:    u,
	}
}

var _ IReactionService = (*reactionService)(nil)

func (rs *reactionService) SetPostReaction(reaction string, postID, userID int) error {
	var isLike bool
	switch reaction {
	case "like":
		isLike = true
	case "dislike":
		isLike = false
	default:
		return entity.ErrInvalidURLPath
	}

	// Build up notification
	userTo, err := rs.postService.GetAuthorID(postID)
	if err != nil {
		return err
	}

	notificaiton := entity.Notification{
		SourceID: postID,
		UserFrom: userID,
		UserTo:   userTo,
	}

	if isLike {
		notificaiton.Type = entity.POST_LIKE
	} else {
		notificaiton.Type = entity.POST_DISLIKE
	}

	// Check if reaction by user exists in table, if so it would return reaction left by that user.
	// If not set new reaction and notify post's author
	isLikeDB, err := rs.reactsRepo.ExistsPostReaction(userID, postID)
	if err != nil {
		if errors.Is(err, entity.ErrNoRecord) {
			err := rs.userService.SendNotification(notificaiton)
			if err != nil {
				return err
			}
			return rs.reactsRepo.AddPostReaction(isLike, postID, userID)
		}
		return err
	}

	err = rs.reactsRepo.DeletePostReaction(postID, userID)
	if err != nil {
		return err
	}

	// If new reaction is not the same as was in table just add new reaction to the table
	// (replacing old reaction with new one) and notify post's author (by deleting old reaction
	// notification and sending new one, basically - replacing old with new one)
	//
	// If it is the same, no change occurs (only deleted old reaction)
	if isLike != isLikeDB {
		var deleteType string
		if isLikeDB {
			deleteType = entity.POST_LIKE
		} else {
			deleteType = entity.POST_DISLIKE
		}
		notificaitonID, err := rs.userService.FindNotification(deleteType, notificaiton.UserFrom, notificaiton.UserTo)
		if err != nil {
			return err
		}
		err = rs.userService.DeleteNotification(notificaitonID)
		if err != nil {
			return err
		}

		err = rs.userService.SendNotification(notificaiton)
		if err != nil {
			return err
		}
		return rs.reactsRepo.AddPostReaction(isLike, postID, userID)
	}

	notificaitonID, err := rs.userService.FindNotification(notificaiton.Type, notificaiton.UserFrom, notificaiton.UserTo)
	if err != nil {
		return err
	}
	err = rs.userService.DeleteNotification(notificaitonID)
	if err != nil {
		return err
	}

	return nil
}

// Same principle to reactions handling in posts
func (rs *reactionService) SetCommentReaction(reaction string, commentID, postID, userID int) error {
	var isLike bool
	switch reaction {
	case "like":
		isLike = true
	case "dislike":
		isLike = false
	default:
		return entity.ErrInvalidURLPath
	}

	userTo, err := rs.commentService.GetAuthorID(commentID)
	if err != nil {
		return err
	}

	notificaiton := entity.Notification{
		SourceID: postID,
		UserFrom: userID,
		UserTo:   userTo,
	}

	if isLike {
		notificaiton.Type = entity.COMMENT_LIKE
	} else {
		notificaiton.Type = entity.COMMENT_DISLIKE
	}

	isLikeDB, err := rs.reactsRepo.ExistsCommentReaction(userID, commentID)
	if err != nil {
		if errors.Is(err, entity.ErrNoRecord) {
			err := rs.userService.SendNotification(notificaiton)
			if err != nil {
				return err
			}
			return rs.reactsRepo.AddCommentReaction(isLike, commentID, userID)
		}
		return err
	}

	err = rs.reactsRepo.DeleteCommentReaction(commentID, userID)
	if err != nil {
		return err
	}

	if isLike != isLikeDB {
		var deleteType string
		if isLikeDB {
			deleteType = entity.POST_LIKE
		} else {
			deleteType = entity.POST_DISLIKE
		}
		notificaitonID, err := rs.userService.FindNotification(deleteType, notificaiton.UserFrom, notificaiton.UserTo)
		if err != nil {
			return err
		}
		err = rs.userService.DeleteNotification(notificaitonID)
		if err != nil {
			return err
		}

		err = rs.userService.SendNotification(notificaiton)
		if err != nil {
			return err
		}
		return rs.reactsRepo.AddCommentReaction(isLike, commentID, userID)
	}

	notificaitonID, err := rs.userService.FindNotification(notificaiton.Type, notificaiton.UserFrom, notificaiton.UserTo)
	if err != nil {
		return err
	}
	err = rs.userService.DeleteNotification(notificaitonID)
	if err != nil {
		return err
	}

	return nil
}
