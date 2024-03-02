package post

import (
	"errors"
	"forum/internal/entity"
	repo "forum/internal/repository/post"
	service "forum/internal/service/post"
	"strconv"
)

type PostServiceMock struct {
	pr repo.IPostRepository
}

func NewPostServiceMock(r repo.IPostRepository) *PostServiceMock {
	return &PostServiceMock{
		pr: r,
	}
}

var _ service.IPostService = (*PostServiceMock)(nil)

func (ps *PostServiceMock) SavePost(p *entity.PostCreateForm, userID int, tags []string) (int, error) {
	if !service.IsRightPost(p) {
		return 0, entity.ErrInvalidFormData
	}

	var tagIDs []int
	for _, tagIDStr := range tags {
		tagID, _ := strconv.Atoi(tagIDStr) // Don't handle error because we know Id's are valid (checked before)
		tagIDs = append(tagIDs, tagID)
	}
	return ps.pr.SavePost(entity.PostCreateForm{}, 0, tagIDs)
}

func (ps *PostServiceMock) GetPost(postID int) (entity.PostView, error) {
	postEntity, err := ps.pr.GetPost(postID)
	if errors.Is(err, entity.ErrNoRecord) {
		return entity.PostView{}, entity.ErrInvalidPostID
	}

	return entity.PostView{
		ID:          postEntity.ID,
		Title:       postEntity.Title,
		Content:     postEntity.Content,
		Username:    postEntity.Username,
		Likes:       postEntity.Likes,
		Dislikes:    postEntity.Dislikes,
		PostTags:    service.ConvertToStrArr(postEntity.PostTags),
		CommentsLen: postEntity.CommentsLen,
	}, nil
}

func (ps *PostServiceMock) GetAllPosts() (*[]entity.PostView, error) {
	posts, _ := ps.pr.GetAllPosts()
	return service.ConvertEntitiesToViews(posts)
}

func (ps *PostServiceMock) GetAllPostsByTagId(tagID int) (*[]entity.PostView, error) {
	posts, _ := ps.pr.GetAllPostsByTagId(tagID)
	return service.ConvertEntitiesToViews(posts)
}

func (ps *PostServiceMock) GetAllPostsByUserId(userID int) (*[]entity.PostView, error) {
	posts, _ := ps.pr.GetAllPostsByUserID(userID)
	return service.ConvertEntitiesToViews(posts)
}

func (ps *PostServiceMock) GetAllPostsByUserReaction(userID int) (*[]entity.PostView, error) {
	posts, _ := ps.pr.GetAllPostsByUserReaction(userID)
	return service.ConvertEntitiesToViews(posts)
}

func (ps *PostServiceMock) ExistsPost(postID int) (bool, error) {
	return ps.pr.ExistsPost(postID)
}
