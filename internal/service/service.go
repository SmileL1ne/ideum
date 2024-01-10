package service

import "forum/internal/entity"

type (
	// TODO: Add new methods related to post manipulation (delete, update post)
	PostService interface {
		SavePost(entity.Post) error
		// GetAllPosts() ([]entity.Post, error)
	}

	// TODO: Same as above
)
