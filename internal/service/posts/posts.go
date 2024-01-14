package posts

import (
	"forum/internal/entity"
	"forum/internal/repository/posts"
)

// TODO: Add new methods related to post manipulation (delete, update post)
type PostService interface {
	GetPost(postId int) (entity.Post, error)
	GetAllPosts() ([]entity.Post, error)
	SavePost(entity.Post) error
}

type postService struct {
	repo posts.PostRepository
}

// Constructor for posts service
func NewPostsService(r posts.PostRepository) *postService {
	return &postService{
		repo: r,
	}
}

func (ps *postService) GetPost(postId int) (entity.Post, error) {
	return entity.Post{}, nil
}

func (uc *postService) GetAllPosts() ([]entity.Post, error) {
	return nil, nil
}

// TODO: Implement this method
func (ps *postService) SavePost(p entity.Post) error {
	err := ps.repo.SavePost(p)
	if err != nil {
		panic(err)
	}
	return nil
}
