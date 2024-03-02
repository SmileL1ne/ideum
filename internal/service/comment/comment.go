package comment

import (
	"forum/internal/entity"
	"forum/internal/repository/comment"
)

type ICommentService interface {
	SaveComment(*entity.CommentCreateForm, int, int) error
	GetAllCommentsForPost(int) (*[]entity.CommentView, error)
	ExistsComment(int) (bool, error)
}

type commentService struct {
	commentRepo comment.ICommentRepository
}

// This ensures that commentService struct implements ICommentService interface
var _ ICommentService = (*commentService)(nil)

func NewCommentService(r comment.ICommentRepository) *commentService {
	return &commentService{
		commentRepo: r,
	}
}

func (cs *commentService) SaveComment(c *entity.CommentCreateForm, postID int, userID int) error {
	if !IsRightComment(c) {
		return entity.ErrInvalidFormData
	}

	err := cs.commentRepo.SaveComment(*c, postID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (cs *commentService) GetAllCommentsForPost(postID int) (*[]entity.CommentView, error) {
	comments, err := cs.commentRepo.GetAllCommentsForPost(postID)
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

func (cs *commentService) ExistsComment(commentID int) (bool, error) {
	return cs.commentRepo.ExistsComment(commentID)
}
