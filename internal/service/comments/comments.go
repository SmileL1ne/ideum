package comments

import (
	"forum/internal/entity"
	"forum/internal/repository/comments"
	"strconv"
)

type ICommentService interface {
	SaveComment(*entity.CommentCreateForm, string) error
	GetAllCommentsForPost(string) (*[]entity.CommentView, error)
}

type commentService struct {
	commentRepo comments.ICommentRepository
}

// This ensures that commentService struct implements ICommentService interface
var _ ICommentService = (*commentService)(nil)

func NewCommentService(r comments.ICommentRepository) *commentService {
	return &commentService{
		commentRepo: r,
	}
}

func (cs *commentService) SaveComment(c *entity.CommentCreateForm, postIdStr string) error {
	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		return entity.ErrInvalidPostId
	}

	if !isRightComment(c) {
		return entity.ErrInvalidFormData
	}

	err = cs.commentRepo.SaveComment(*c, postId)
	if err != nil {
		return err
	}

	return nil
}

func (cs *commentService) GetAllCommentsForPost(postIdStr string) (*[]entity.CommentView, error) {
	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		return nil, entity.ErrInvalidPostId
	}

	comments, err := cs.commentRepo.GetAllCommentsForPost(postId)

	// Convert received CommentEntity's to CommentView's
	var cViews []entity.CommentView
	for _, c := range *comments {
		comment := entity.CommentView{
			Content:   c.Content,
			CreatedAt: c.CreatedAt,
		}
		cViews = append(cViews, comment)
	}

	return &cViews, nil
}
