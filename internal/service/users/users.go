package users

import (
	"errors"
	"forum/internal/entity"
	"forum/internal/repository/users"
	"forum/internal/validator"

	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	SaveUser(*entity.UserSignupForm) (int, error)
	Authenticate(*entity.UserLoginForm) (int, error)
}

type userService struct {
	userRepo users.IUserRepository
}

func NewUserService(u users.IUserRepository) *userService {
	return &userService{
		userRepo: u,
	}
}

var _ IUserService = (*userService)(nil)

func (us *userService) SaveUser(u *entity.UserSignupForm) (int, error) {
	if !isRightSignUp(u) {
		return 0, entity.ErrInvalidFormData
	}

	id, err := us.userRepo.SaveUser(*u)
	if err != nil {
		if errors.Is(err, entity.ErrDuplicateEmail) {
			u.AddFieldError("email", "Email address is already in use")
			return 0, entity.ErrInvalidFormData
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (us *userService) Authenticate(u *entity.UserLoginForm) (int, error) {
	if !isRightLogin(u) {
		return 0, entity.ErrInvalidFormData
	}

	var userFromDB entity.UserEntity
	var err error

	if validator.Matches(u.Identifier, EmailRX) {
		userFromDB, err = us.userRepo.GetUserByEmail(u.Identifier)
	} else {
		userFromDB, err = us.userRepo.GetUserByUsername(u.Identifier)
	}
	if err != nil {
		if errors.Is(err, entity.ErrInvalidCredentials) {
			u.AddNonFieldError("Email or password is incorrect")
			return 0, entity.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(userFromDB.Password), []byte(u.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			u.AddNonFieldError("Email or password is incorrect")
			return 0, entity.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return userFromDB.Id, nil
}
