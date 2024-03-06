package post

import (
	"errors"
	"forum/internal/entity"
	repo "forum/internal/repository/post"
	"forum/internal/service/image"
	service "forum/internal/service/post"
	"forum/internal/service/tag"
	"strconv"
)

type PostServiceMock struct {
	imgService image.IImageService
	tagService tag.ITagService
	pr         repo.IPostRepository
}

func NewPostServiceMock(r repo.IPostRepository, is image.IImageService, ts tag.ITagService) *PostServiceMock {
	return &PostServiceMock{
		imgService: is,
		tagService: ts,
		pr:         r,
	}
}

var _ service.IPostService = (*PostServiceMock)(nil)

func (ps *PostServiceMock) SavePost(p entity.PostCreateForm) (int, error) {

	var tagIDs []int
	for _, tagIDStr := range p.Tags {
		tagID, _ := strconv.Atoi(tagIDStr) // Don't handle error because we know Id's are valid (checked before)
		tagIDs = append(tagIDs, tagID)
	}
	return ps.pr.SavePost(entity.PostCreateForm{}, tagIDs)
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

// TODO: Add mock checks here
func (ps *PostServiceMock) CheckPostAttrs(p *entity.PostCreateForm, withImage bool) (bool, error) {
	if !service.IsRightPost(p, withImage) {
		return false, entity.ErrInvalidFormData
	}

	areTagsExist, err := ps.tagService.AreTagsExist(p.Tags)
	if !areTagsExist || err != nil {
		return false, err
	}

	return false, nil
}
