package usecase

import "forum/internal/entity"

type (
	// TODO: Add new methods related to post manipulation (delete, update post)
	Post interface {
		MakeNewPost(entity.Post) error
		// GetAllPosts() ([]entity.Post, error)
	}

	// TODO: Same as above
	PostRepo interface {
		SavePost(entity.Post) error
	}
)
