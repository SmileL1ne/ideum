package entity

import (
	"forum/internal/validator"
	"time"
)

// UserEntity is returned by repos (not pointer, because data is read only)
type UserEntity struct {
	Id        int
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
}

// UserSignupForm is accepted by services as pointers to save form validation
// errors purpose only. It accepted as copy by repos because it is read only
// at that stage
type UserSignupForm struct {
	Username string
	Email    string
	Password string
	validator.Validator
}

// UserLoginForm is accepted by services as pointers, by repos as copies for the
// same purposes as UserSignupForm
type UserLoginForm struct {
	Identifier string
	Password   string
	validator.Validator
}
