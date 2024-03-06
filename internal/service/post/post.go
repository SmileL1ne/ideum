package post

import (
	"errors"
	"forum/internal/entity"
	"forum/internal/repository/post"
	"forum/internal/service/image"
	"forum/internal/service/tag"
	"strconv"
)

type IPostService interface {
	SavePost(entity.PostCreateForm) (int, error)
	GetPost(int) (entity.PostView, error)
	GetAllPosts() (*[]entity.PostView, error)
	GetAllPostsByTagId(int) (*[]entity.PostView, error)
	GetAllPostsByUserId(int) (*[]entity.PostView, error)
	GetAllPostsByUserReaction(int) (*[]entity.PostView, error)
	ExistsPost(int) (bool, error)
	CheckPostAttrs(*entity.PostCreateForm, bool) (bool, error)
}

type postService struct {
	imgService image.IImageService
	tagService tag.ITagService
	postRepo   post.IPostRepository
}

// Constructor for post service
func NewPostsService(r post.IPostRepository, is image.IImageService, ts tag.ITagService) *postService {
	return &postService{
		imgService: is,
		tagService: ts,
		postRepo:   r,
	}
}

var _ IPostService = (*postService)(nil)

func (ps *postService) SavePost(p entity.PostCreateForm) (int, error) {
	var tagIDs []int
	for _, tagIDStr := range p.Tags {
		tagID, _ := strconv.Atoi(tagIDStr) // Don't handle error because we know Ids are valid (checked before)
		tagIDs = append(tagIDs, tagID)
	}

	id, err := ps.postRepo.SavePost(p, tagIDs)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (ps *postService) GetPost(postId int) (entity.PostView, error) {
	post, err := ps.postRepo.GetPost(postId)
	if err != nil {
		if errors.Is(err, entity.ErrNoRecord) {
			return entity.PostView{}, entity.ErrInvalidPostID
		}
		return entity.PostView{}, err
	}

	imgName, err := ps.imgService.Get(postId)
	if err != nil {
		return entity.PostView{}, err
	}

	tags := ConvertToStrArr(post.PostTags)
	pView := entity.PostView{
		ID:          post.ID,
		Title:       post.Title,
		Content:     post.Content,
		CreatedAt:   post.CreatedAt,
		Username:    post.Username,
		Likes:       post.Likes,
		Dislikes:    post.Dislikes,
		CommentsLen: post.CommentsLen,
		PostTags:    tags,
		ImageName:   imgName,
	}

	return pView, nil
}

func (ps *postService) GetAllPosts() (*[]entity.PostView, error) {
	posts, err := ps.postRepo.GetAllPosts()
	if err != nil {
		return nil, err
	}

	return ConvertEntitiesToViews(posts)
}

func (ps *postService) GetAllPostsByTagId(tagID int) (*[]entity.PostView, error) {
	posts, err := ps.postRepo.GetAllPostsByTagId(tagID)
	if err != nil {
		return nil, err
	}

	return ConvertEntitiesToViews(posts)
}

func (ps *postService) GetAllPostsByUserId(userID int) (*[]entity.PostView, error) {
	posts, err := ps.postRepo.GetAllPostsByUserID(userID)
	if err != nil {
		return nil, err
	}

	return ConvertEntitiesToViews(posts)
}

func (ps *postService) GetAllPostsByUserReaction(userID int) (*[]entity.PostView, error) {
	posts, err := ps.postRepo.GetAllPostsByUserReaction(userID)
	if err != nil {
		return nil, err
	}

	return ConvertEntitiesToViews(posts)
}

func (ps *postService) ExistsPost(postID int) (bool, error) {
	return ps.postRepo.ExistsPost(postID)
}

func (ps *postService) CheckPostAttrs(p *entity.PostCreateForm, withImage bool) (bool, error) {
	if !IsRightPost(p, withImage) {
		return false, nil
	}

	areTagsExist, err := ps.tagService.AreTagsExist(p.Tags)
	if !areTagsExist || err != nil {
		return false, err
	}

	return true, nil
}
