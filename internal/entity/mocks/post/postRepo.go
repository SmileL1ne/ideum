package post

import (
	"forum/internal/entity"
	"forum/internal/repository/post"
	"time"
)

var mockPost = entity.PostEntity{
	ID:          1,
	Title:       "Shine bright",
	Content:     "You can do it!",
	CreatedAt:   time.Date(2003, time.July, 6, 0, 0, 0, 0, time.Local),
	UserID:      7,
	Username:    "mustik",
	Likes:       5,
	Dislikes:    5,
	PostTags:    "Art, Games",
	CommentsLen: 2,
}

type PostRepoMock struct{}

func NewPostRepoMock() *PostRepoMock {
	return &PostRepoMock{}
}

var _ post.IPostRepository = (*PostRepoMock)(nil)

func (r *PostRepoMock) Insert(p entity.PostCreateForm, tagIDs []int) (int, error) {
	return mockPost.ID, nil
}

func (r *PostRepoMock) Get(postID int) (entity.PostEntity, error) {
	switch postID {
	case mockPost.ID:
		return mockPost, nil
	default:
		return mockPost, entity.ErrNoRecord
	}
}

func (r *PostRepoMock) GetAll() (*[]entity.PostEntity, error) {
	return &[]entity.PostEntity{mockPost}, nil
}

func (r *PostRepoMock) GetAllByTagId(tagID int) (*[]entity.PostEntity, error) {
	return &[]entity.PostEntity{mockPost}, nil
}

func (r *PostRepoMock) GetAllByUserID(userID int) (*[]entity.PostEntity, error) {
	return &[]entity.PostEntity{mockPost}, nil
}

func (r *PostRepoMock) GetAllByUserReaction(userID int) (*[]entity.PostEntity, error) {
	return &[]entity.PostEntity{mockPost}, nil
}

func (r *PostRepoMock) Exists(postID int) (bool, error) {
	if postID != mockPost.ID {
		return false, nil
	}
	return true, nil
}

func (r *PostRepoMock) Delete(postID int) error {
	return nil
}
