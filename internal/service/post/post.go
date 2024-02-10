package post

import (
	"errors"
	"forum/internal/entity"
	"forum/internal/repository/post"
	"strconv"
)

/*
TODO:
- Add new methods related to post manipulation (delete, update post)
- Add data validation to SavePost method
- Add postID validation in GetPost method
*/
type IPostService interface {
	SavePost(*entity.PostCreateForm, int, []string) (int, error)
	GetPost(int) (entity.PostView, error)
	GetAllPosts() (*[]entity.PostView, error)
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
	if !isRightPost(p) {
		return 0, entity.ErrInvalidFormData
	}
	var tagIDs []int
	for _, tagIDStr := range tags {
		tagID, err := strconv.Atoi(tagIDStr)
		if err != nil {
			return 0, entity.ErrInvalidFormData
		}
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
			return entity.PostView{}, err
		}
		return entity.PostView{}, err
	}

	pView := entity.PostView{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		CreatedAt: post.CreatedAt,
		Username:  post.Username,
		Likes:     post.Likes,
		Dislikes:  post.Dislikes,
	}

	return pView, nil
}

func (uc *postService) GetAllPosts() (*[]entity.PostView, error) {
	posts, err := uc.postRepo.GetAllPosts()
	if err != nil {
		return nil, err
	}

	// Convert received PostEntity's to PostView's
	var pViews []entity.PostView
	for _, p := range *posts {
		post := entity.PostView{
			ID:        p.ID,
			Title:     p.Title,
			Content:   p.Content,
			CreatedAt: p.CreatedAt,
			Username:  p.Username,
			Likes:     p.Likes,
			Dislikes:  p.Dislikes,
		}
		pViews = append(pViews, post)
	}

	return &pViews, nil
}
