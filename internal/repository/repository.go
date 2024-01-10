package repository

import "forum/internal/entity"

type (

	// TODO: Add new methods related to post manipulation (delete, update post)
	PostRepository interface {
		SavePost(p entity.Post) error
		// RetrieveAllPosts() ([]entity.Post, error)
	}
)
