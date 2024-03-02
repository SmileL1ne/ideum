package post

import (
	"errors"
	"forum/internal/entity"
	"forum/internal/repository/post"
	"strconv"
)

type IPostService interface {
	SavePost(*entity.PostCreateForm, int, []string) (int, error)
	GetPost(int) (entity.PostView, error)
	GetAllPosts() (*[]entity.PostView, error)
	GetAllPostsByTagId(int) (*[]entity.PostView, error)
	GetAllPostsByUserId(int) (*[]entity.PostView, error)
	GetAllPostsByUserReaction(int) (*[]entity.PostView, error)
	ExistsPost(int) (bool, error)
}

type postService struct {
	postRepo post.IPostRepository
}

// Constructor for post service
func NewPostsService(r post.IPostRepository) *postService {
	return &postService{
		postRepo: r,
	}
}

var _ IPostService = (*postService)(nil)

func (ps *postService) SavePost(p *entity.PostCreateForm, userID int, tags []string) (int, error) {
	if !IsRightPost(p) {
		return 0, entity.ErrInvalidFormData
	}

	var tagIDs []int
	for _, tagIDStr := range tags {
		tagID, _ := strconv.Atoi(tagIDStr) // Don't handle error because we know Id's are valid (checked before)
		tagIDs = append(tagIDs, tagID)
	}

	id, err := ps.postRepo.SavePost(*p, userID, tagIDs)
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
	}

	return pView, nil
}

func (pc *postService) GetAllPosts() (*[]entity.PostView, error) {
	posts, err := pc.postRepo.GetAllPosts()
	if err != nil {
		return nil, err
	}

	return ConvertEntitiesToViews(posts)
}

func (pc *postService) GetAllPostsByTagId(tagID int) (*[]entity.PostView, error) {
	posts, err := pc.postRepo.GetAllPostsByTagId(tagID)
	if err != nil {
		return nil, err
	}

	return ConvertEntitiesToViews(posts)
}

func (pc *postService) GetAllPostsByUserId(userID int) (*[]entity.PostView, error) {
	posts, err := pc.postRepo.GetAllPostsByUserID(userID)
	if err != nil {
		return nil, err
	}

	return ConvertEntitiesToViews(posts)
}

func (pc *postService) GetAllPostsByUserReaction(userID int) (*[]entity.PostView, error) {
	posts, err := pc.postRepo.GetAllPostsByUserReaction(userID)
	if err != nil {
		return nil, err
	}

	return ConvertEntitiesToViews(posts)
}

func (pc *postService) ExistsPost(postID int) (bool, error) {
	return pc.postRepo.ExistsPost(postID)
}
