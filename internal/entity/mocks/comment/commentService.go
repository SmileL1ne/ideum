package comment

import (
	"forum/internal/entity"
	repo "forum/internal/repository/comment"
	service "forum/internal/service/comment"
)

type CommentServiceMock struct {
	cr repo.ICommentRepository
}

func NewCommentServiceMock(r repo.ICommentRepository) *CommentServiceMock {
	return &CommentServiceMock{
		cr: r,
	}
}

var _ service.ICommentService = (*CommentServiceMock)(nil)

func (cs *CommentServiceMock) SaveComment(c *entity.CommentCreateForm, postID int, userID int) error {
	if !service.IsRightComment(c) {
		return entity.ErrInvalidFormData
	}

	return cs.cr.Insert(*c, postID, userID)
}

func (cs *CommentServiceMock) GetAllCommentsForPost(postID int) (*[]entity.CommentView, error) {
	comments, err := cs.cr.GetAllForPost(postID)
	if err != nil {
		return nil, err
	}

	return service.ConvertEntitiesToViews(comments)
}

func (cs *CommentServiceMock) GetAllUserCommentsForPost(userID, postID int) (*[]entity.CommentView, error) {
	comments, err := cs.cr.GetAllUserCommentsForPost(userID, postID)
	if err != nil {
		return nil, err
	}

	return service.ConvertEntitiesToViews(comments)
}

func (cs *CommentServiceMock) ExistsComment(commentID int) (bool, error) {
	return cs.cr.Exists(commentID)
}
