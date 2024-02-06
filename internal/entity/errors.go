package entity

import "errors"

var (
	// Common errors
	ErrNoRecord        = errors.New("entity: no matching row found")
	ErrInvalidFormData = errors.New("entity: some form data is invalid")
)

var (
	// Post related errors
	ErrInvalidPostId = errors.New("entity: invalid post id")
)

var (
	// User related errors
	ErrDuplicateEmail     = errors.New("entity: duplicate email")
	ErrInvalidCredentials = errors.New("entity: invalid credentials")
)
