package entity

import (
	"forum/internal/validator"
	"time"
)

// CommentEntity is returned by repositories, storing all comment related data
// for business logic
type CommentEntity struct {
	Id        int
	Content   string
	CreatedAt time.Time
	PostID    int
	/*
		TODO:
		- Add 'Username' field to use when retrieving all comments left by particular user
	*/
	//
}

// CommentView is returned by services, storing all comment related data that
// will be outputed to end user
type CommentView struct {
	Content   string
	CreatedAt time.Time
}

// CommentCreateForm is accepted by services by pointer only for form error messages
// handling, so they ar written in Validator's FieldErrors or NonFieldErrors fields
type CommentCreateForm struct {
	Content string
	validator.Validator
}
