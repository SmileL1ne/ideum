package entity

import (
	"forum/internal/validator"
	"time"
)

// Comment
type CommentEntity struct {
	Id        int
	Content   string
	CreatedAt time.Time
	PostId    int
	// UserId int - add in future
}

type CommentView struct {
	Content   string
	CreatedAt time.Time
}

type CommentCreateForm struct {
	Content string
	validator.Validator
}
