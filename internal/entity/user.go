package entity

import (
	"forum/internal/service/validator"
	"time"
)

type UserEntity struct {
	Id        int
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
}

type UserSignupForm struct {
	Username string
	Email    string
	Password string
	validator.Validator
}
