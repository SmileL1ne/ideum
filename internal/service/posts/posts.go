package posts

import (
	"errors"
	"forum/internal/entity"
	"forum/internal/repository/posts"
	"net/http"
	"strconv"
)

/*
TODO:
- Add new methods related to post manipulation (delete, update post)
- Add data validation to SavePost method
- Add postID validation in GetPost method
*/
type PostService interface {
	SavePost(*entity.PostCreateForm) (int, int, error)
	GetPost(postId string) (entity.PostEntity, int, error)
	GetAllPosts() ([]entity.PostEntity, error)
}

type postService struct {
	postRepo posts.PostRepository
}

// Constructor for post service
func NewPostsService(r posts.PostRepository) *postService {
	return &postService{
		postRepo: r,
	}
}

func (ps *postService) SavePost(p *entity.PostCreateForm) (int, int, error) {
	if !isRightPost(p) {
		return 0, http.StatusUnprocessableEntity, nil
	}

	id, err := ps.postRepo.SavePost(*p)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	return id, http.StatusOK, nil
}

func (ps *postService) GetPost(postId string) (entity.PostEntity, int, error) {

	// PostId validation here

	id, err := strconv.Atoi(postId)
	if err != nil {
		return entity.PostEntity{}, http.StatusBadRequest, err
	}

	post, err := ps.postRepo.GetPost(id)
	if err != nil {
		if errors.Is(err, entity.ErrNoRecord) {
			return entity.PostEntity{}, http.StatusNotFound, err
		}
		return entity.PostEntity{}, http.StatusInternalServerError, err
	}

	return post, http.StatusOK, nil
}

func (uc *postService) GetAllPosts() ([]entity.PostEntity, error) {
	return uc.postRepo.GetAllPosts()
}
