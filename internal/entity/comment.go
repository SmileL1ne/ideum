package entity

import (
	"forum/internal/validator"
	"time"
)

// CommentEntity is returned by repositories, storing all comment related data
// for business logic
type CommentEntity struct {
	ID        int
	Content   string
	CreatedAt time.Time
	PostID    int
	Username  string
	Likes     int
	Dislikes  int
}

// CommentView is returned by services, storing all comment related data that
// will be outputed to end user
type CommentView struct {
	ID        int
	Username  string
	Content   string
	CreatedAt time.Time
	PostID    int
	Likes     int
	Dislikes  int
}

// CommentCreateForm is accepted by services by pointer only for form error messages
// handling, so they ar written in Validator's FieldErrors or NonFieldErrors fields
type CommentCreateForm struct {
	Content string
	validator.Validator
}
