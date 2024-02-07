package posts

import (
	"errors"
	"forum/internal/entity"
	"forum/internal/repository/posts"
	"strconv"
)

/*
TODO:
- Add new methods related to post manipulation (delete, update post)
- Add data validation to SavePost method
- Add postID validation in GetPost method
*/
type IPostService interface {
	SavePost(*entity.PostCreateForm) (int, error)
	GetPost(string) (entity.PostView, error)
	GetAllPosts() (*[]entity.PostView, error)
}

type postService struct {
	postRepo posts.IPostRepository
}

// Constructor for post service
func NewPostsService(r posts.IPostRepository) *postService {
	return &postService{
		postRepo: r,
	}
}

var _ IPostService = (*postService)(nil)

func (ps *postService) SavePost(p *entity.PostCreateForm) (int, error) {
	if !isRightPost(p) {
		return 0, entity.ErrInvalidFormData
	}

	id, err := ps.postRepo.SavePost(*p)
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
		Id:        post.Id,
		Title:     post.Title,
		Content:   post.Content,
		CreatedAt: post.Created,
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
			Id:        p.Id,
			Title:     p.Title,
			Content:   p.Content,
			CreatedAt: p.Created,
		}
		pViews = append(pViews, post)
	}

	return &pViews, nil
}
