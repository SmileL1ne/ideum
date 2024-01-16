package entity

import "errors"

var (
	ErrNoRecord = errors.New("entity: no matching row found")
)
