package entity

import "errors"

// Common errors
var (
	ErrNoRecord        = errors.New("entity: no matching row found")
	ErrInvalidFormData = errors.New("entity: some form data is invalid")
)

// Post related errors
var (
	ErrInvalidPostId = errors.New("entity: invalid post id")
)

// Comment related errors
var (
	ErrInvalidCommentId = errors.New("entity: invalid comment id")
)

// User related errors
var (
	ErrDuplicateEmail     = errors.New("entity: duplicate email")
	ErrDuplicateUsername  = errors.New("entity: dupliate username")
	ErrInvalidCredentials = errors.New("entity: invalid credentials")
)
