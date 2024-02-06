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
	GetPost(postId string) (entity.PostEntity, error)
	GetAllPosts() ([]entity.PostEntity, error)
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

func (ps *postService) GetPost(postId string) (entity.PostEntity, error) {
	id, err := strconv.Atoi(postId)
	if err != nil {
		return entity.PostEntity{}, entity.ErrInvalidPostId
	}

	post, err := ps.postRepo.GetPost(id)
	if err != nil {
		if errors.Is(err, entity.ErrNoRecord) {
			return entity.PostEntity{}, err
		}
		return entity.PostEntity{}, err
	}

	return post, nil
}

func (uc *postService) GetAllPosts() ([]entity.PostEntity, error) {
	return uc.postRepo.GetAllPosts()
}
