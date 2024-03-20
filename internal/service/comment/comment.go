package comment

import (
	"forum/internal/entity"
	"forum/internal/repository/comment"
	"forum/internal/service/user"
)

type ICommentService interface {
	SaveComment(*entity.CommentCreateForm, int, int) error
	GetAllCommentsForPost(int) (*[]entity.CommentView, error)
	GetAllUserCommentsForPost(userID, postID int) (*[]entity.CommentView, error)
	ExistsComment(int) (bool, error)
	DeleteComment(commentID, userID int) error
	DeleteCommentPrivileged(commentID int, userRole string) error
}

type commentService struct {
	commentRepo comment.ICommentRepository
	userService user.IUserService
}

// This ensures that commentService struct implements ICommentService interface
var _ ICommentService = (*commentService)(nil)

func NewCommentService(r comment.ICommentRepository, us user.IUserService) *commentService {
	return &commentService{
		commentRepo: r,
		userService: us,
	}
}

func (cs *commentService) SaveComment(c *entity.CommentCreateForm, postID int, userID int) error {
	if !IsRightComment(c) {
		return entity.ErrInvalidFormData
	}

	err := cs.commentRepo.Insert(*c, postID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (cs *commentService) GetAllCommentsForPost(postID int) (*[]entity.CommentView, error) {
	comments, err := cs.commentRepo.GetAllForPost(postID)
	if err != nil {
		return nil, err
	}

	return ConvertEntitiesToViews(comments)
}

func (cs *commentService) GetAllUserCommentsForPost(userID, postID int) (*[]entity.CommentView, error) {
	comments, err := cs.commentRepo.GetAllUserCommentsForPost(userID, postID)
	if err != nil {
		return nil, err
	}

	return ConvertEntitiesToViews(comments)
}

func (cs *commentService) ExistsComment(commentID int) (bool, error) {
	return cs.commentRepo.Exists(commentID)
}

func (cs *commentService) DeleteComment(commentID, userID int) error {
	exists, err := cs.commentRepo.Exists(commentID)
	if err != nil {
		return err
	}
	if !exists {
		return entity.ErrCommentNotFound
	}

	return cs.commentRepo.Delete(commentID, userID)
}

func (cs *commentService) DeleteCommentPrivileged(commentID int, userRole string) error {
	exists, err := cs.commentRepo.Exists(commentID)
	if err != nil {
		return err
	}
	if !exists {
		return entity.ErrCommentNotFound
	}

	authorID, err := cs.commentRepo.GetAuthorID(commentID)
	if err != nil {
		return err
	}

	err = cs.commentRepo.DeleteByPrivileged(commentID)
	if err != nil {
		return err
	}

	notificaiton := entity.Notification{
		Type:     entity.DELETE_COMMENT,
		UserFrom: authorID,
	}
	if userRole == entity.MODERATOR {
		notificaiton.Content = ". Reason: obscene"
	}

	err = cs.userService.SendNotification(notificaiton)
	if err != nil {
		return err
	}

	return nil
}
