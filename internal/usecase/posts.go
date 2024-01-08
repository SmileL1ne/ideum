package usecase

import "forum/internal/entity"

type PostsUseCase struct {
	repo PostRepo
}

func New(r PostRepo) *PostsUseCase {
	return &PostsUseCase{
		repo: r,
	}
}

// TODO: Implement this method
func (uc *PostsUseCase) MakeNewPost(p entity.Post) error {
	err := uc.repo.SavePost(p)
	if err != nil {
		panic(err)
	}
	return nil
}

// func (uc *PostsUseCase) GetAllPosts() ([]entity.Post, error)
