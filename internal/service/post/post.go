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
	SavePost(*entity.PostCreateForm, int) (int, error)
	GetPost(string) (entity.PostView, error)
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

func (ps *postService) SavePost(p *entity.PostCreateForm, userID int) (int, error) {
	if !isRightPost(p) {
		return 0, entity.ErrInvalidFormData
	}

	id, err := ps.postRepo.SavePost(*p, userID)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (ps *postService) GetPost(postId string) (entity.PostView, error) {
	id, err := strconv.Atoi(postId)
	if err != nil {
		return entity.PostView{}, entity.ErrInvalidPostId
	}

	post, err := ps.postRepo.GetPost(id)
	if err != nil {
		if errors.Is(err, entity.ErrNoRecord) {
			return entity.PostView{}, err
		}
		return entity.PostView{}, err
	}

	pView := entity.PostView{
		Id:        post.ID,
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
			Id:        p.ID,
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
