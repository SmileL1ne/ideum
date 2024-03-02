package comment

import (
	"forum/internal/entity"
	repo "forum/internal/repository/comment"
	service "forum/internal/service/comment"
)

type CommentServiceMock struct {
	cr repo.ICommentRepository
}

func NewTagServiceMock(r repo.ICommentRepository) *CommentServiceMock {
	return &CommentServiceMock{
		cr: r,
	}
}

var _ service.ICommentService = (*CommentServiceMock)(nil)

func (cs *CommentServiceMock) SaveComment(c *entity.CommentCreateForm, postID int, userID int) error {
	if !service.IsRightComment(c) {
		return entity.ErrInvalidFormData
	}

	return cs.cr.SaveComment(*c, postID, userID)
}

func (cs *CommentServiceMock) GetAllCommentsForPost(postID int) (*[]entity.CommentView, error) {
	comments, err := cs.cr.GetAllCommentsForPost(postID)
	if err != nil {
		return nil, err
	}

	// Convert received CommentEntity's to CommentView's
	var cViews []entity.CommentView
	for _, c := range *comments {
		comment := entity.CommentView{
			ID:        c.ID,
			Username:  c.Username,
			Content:   c.Content,
			CreatedAt: c.CreatedAt,
			PostID:    c.PostID,
			Likes:     c.Likes,
			Dislikes:  c.Dislikes,
		}
		cViews = append(cViews, comment)
	}

	return &cViews, nil
}

func (cs *CommentServiceMock) ExistsComment(commentID int) (bool, error) {
	return cs.cr.ExistsComment(commentID)
}
