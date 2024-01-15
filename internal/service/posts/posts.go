package posts

import (
	"forum/internal/entity"
	"forum/internal/repository/posts"
)

// TODO: Add new methods related to post manipulation (delete, update post)
type PostService interface {
	SavePost(entity.Post) (int, error)
	GetPost(postId int) (entity.Post, error)
	GetAllPosts() ([]entity.Post, error)
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

// TODO: Implement this method
func (ps *postService) SavePost(p entity.Post) (int, error) {
	id, err := ps.repo.SavePost(p)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (ps *postService) GetPost(postId int) (entity.Post, error) {
	return entity.Post{}, nil
}

func (uc *postService) GetAllPosts() ([]entity.Post, error) {
	return nil, nil
}
