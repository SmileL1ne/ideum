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

func (r *PostRepoMock) SavePost(p entity.PostCreateForm, userID int, tagIDs []int) (int, error) {
	return 2, nil
}

func (r *PostRepoMock) GetPost(postID int) (entity.PostEntity, error) {
	switch postID {
	case 1:
		return mockPost, nil
	default:
		return entity.PostEntity{}, entity.ErrNoRecord
	}
}

func (r *PostRepoMock) GetAllPosts() (*[]entity.PostEntity, error) {
	return nil, nil
}

func (r *PostRepoMock) GetAllPostsByTagId(tagID int) (*[]entity.PostEntity, error) {
	return nil, nil
}

func (r *PostRepoMock) GetAllPostsByUserID(userID int) (*[]entity.PostEntity, error) {
	return nil, nil
}

func (r *PostRepoMock) GetAllPostsByUserReaction(userID int) (*[]entity.PostEntity, error) {
	return nil, nil
}

func (r *PostRepoMock) ExistsPost(postID int) (bool, error) {
	return true, nil
}