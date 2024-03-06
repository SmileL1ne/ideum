package entity

import (
	"forum/internal/validator"
	"mime/multipart"
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
	ImageName   string
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
	ImageName   string
}

// PostCreateForm is accepted by services and repos. Services accept pointer only for
// form error messages handling, so they are written in Validator's FieldErrors
// or NonFieldErrors fields
type PostCreateForm struct {
	Title      string
	Content    string
	UserID     int
	Tags       []string
	File       multipart.File
	FileHeader *multipart.FileHeader
	ImageName  string
	validator.Validator
}
