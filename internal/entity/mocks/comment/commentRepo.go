package comment

import (
	"forum/internal/entity"
	"forum/internal/repository/comment"
	"time"
)

var mockComment = entity.CommentEntity{
	ID:        1,
	Content:   "hell nah",
	CreatedAt: time.Date(2003, time.July, 6, 0, 0, 0, 0, time.Local),
	PostID:    1,
	Username:  "mustik",
	Likes:     5,
	Dislikes:  5,
}

type CommentRepoMock struct{}

func NewCommentRepoMock() *CommentRepoMock {
	return &CommentRepoMock{}
}

var _ comment.ICommentRepository = (*CommentRepoMock)(nil)

func (r *CommentRepoMock) SaveComment(c entity.CommentCreateForm, postID int, userID int) error {
	return nil
}

func (r *CommentRepoMock) GetAllCommentsForPost(postID int) (*[]entity.CommentEntity, error) {
	return &[]entity.CommentEntity{mockComment}, nil
}

func (r *CommentRepoMock) ExistsComment(commentID int) (bool, error) {
	return true, nil
}
