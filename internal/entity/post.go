package entity

import (
	"forum/internal/validator"
	"time"
)

// TODO: Add database related metatags

type PostEntity struct {
	Id      int
	Title   string
	Content string
	Created time.Time
	// UserID  int // Foreign key
}

type PostCreateForm struct {
	Title   string
	Content string
	validator.Validator
}
