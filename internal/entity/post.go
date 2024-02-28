package entity

import (
	"forum/internal/validator"
	"time"
)

// PostEntity is returned to services from repositories for business logic purposes
type PostEntity struct {
	ID          int
	Title       string
	Content     string
	CreatedAt   time.Time
	UserID      int
	Username    string
	Likes       int
	Dislikes    int
	PostTags    string
	CommentsLen int
}

// PostView is returned to handlers from service and outputed in pages
type PostView struct {
	ID          int
	Title       string
	Content     string
	CreatedAt   time.Time
	Username    string
	Likes       int
	Dislikes    int
	PostTags    []string
	CommentsLen int
}

// PostCreateForm is accepted by services and repos. Services accept pointer only for
// form error messages handling, so they are written in Validator's FieldErrors
// or NonFieldErrors fields
type PostCreateForm struct {
	Title   string
	Content string
	validator.Validator
}
