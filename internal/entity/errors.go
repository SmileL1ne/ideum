package entity

import "errors"

var (
	ErrNoRecord       = errors.New("entity: no matching row found")
	ErrDuplicateEmail = errors.New("entity: duplicate email")
)
