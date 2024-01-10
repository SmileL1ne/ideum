package usecase

import (
	"forum/internal/entity"
	"forum/internal/repository"
)

type postsService struct {
	repo repository.PostRepository
}

func New(r repository.PostRepository) *postsService {
	return &postsService{
		repo: r,
	}
}

// TODO: Implement this method
func (uc *postsService) SavePost(p entity.Post) error {
	err := uc.repo.SavePost(p)
	if err != nil {
		panic(err)
	}
	return nil
}

// func (uc *PostsUseCase) GetAllPosts() ([]entity.Post, error)
