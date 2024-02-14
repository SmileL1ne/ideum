package entity

import "errors"

// Common errors
var (
	ErrNoRecord        = errors.New("entity: no matching row found")
	ErrInvalidFormData = errors.New("entity: some form data is invalid")
)

// Post related errors
var (
	ErrInvalidURLPath = errors.New("entity: invalid url path")
)

// Comment related errors
var ()

// User related errors
var (
	ErrDuplicateEmail     = errors.New("entity: duplicate email")
	ErrDuplicateUsername  = errors.New("entity: dupliate username")
	ErrInvalidCredentials = errors.New("entity: invalid credentials")
	ErrInvalidUserID      = errors.New("entity: non-existent user id")
)
