package comment

import (
	"forum/internal/entity"
	"forum/internal/repository/comment"
	"strconv"
)

type ICommentService interface {
	SaveComment(*entity.CommentCreateForm, string, int) error
	GetAllCommentsForPost(string) (*[]entity.CommentView, error)
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

func (cs *commentService) SaveComment(c *entity.CommentCreateForm, postIDStr string, userID int) error {
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		return entity.ErrInvalidPostId
	}

	if !isRightComment(c) {
		return entity.ErrInvalidFormData
	}

	err = cs.commentRepo.SaveComment(*c, postID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (cs *commentService) GetAllCommentsForPost(postIDStr string) (*[]entity.CommentView, error) {
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		return nil, entity.ErrInvalidPostId
	}

	comments, err := cs.commentRepo.GetAllCommentsForPost(postID)

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
