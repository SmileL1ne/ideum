package entity

import (
	"forum/internal/validator"
	"time"
)

type CommentEntity struct {
	Id        int
	Content   string
	CreatedAt time.Time
	PostId    int
	// UserId int - add in future
}

type CommentCreateForm struct {
	Content string
	validator.Validator
}
