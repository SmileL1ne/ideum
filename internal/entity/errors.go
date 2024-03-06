package entity

import "errors"

// Common errors
var (
	ErrNoRecord        = errors.New("entity: no matching row found")
	ErrInvalidFormData = errors.New("entity: some form data is invalid")
	ErrInvalidPathID   = errors.New("entity: invalid id in request path")
	ErrInvalidURLPath  = errors.New("entity: invalid url path")
)

// Post related errors
var (
	ErrInvalidPostID    = errors.New("entity: invalid post id")
	ErrInvalidImageType = errors.New("entity: invalid image type")
	ErrTooLargeImage    = errors.New("entity: too large image")
	ErrInvalidTags      = errors.New("entity: tags don't exist")
)

// User related errors
var (
	ErrDuplicateEmail     = errors.New("entity: duplicate email")
	ErrDuplicateUsername  = errors.New("entity: dupliate username")
	ErrInvalidCredentials = errors.New("entity: invalid credentials")
	ErrInvalidUserID      = errors.New("entity: non-existent user id")
)
