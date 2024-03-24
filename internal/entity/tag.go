package entity

import (
	"forum/internal/validator"
	"time"
)

// TagEntity is returned from repo and from service
type TagEntity struct {
	ID        int
	Name      string
	CreatedAt time.Time
	validator.Validator
}
