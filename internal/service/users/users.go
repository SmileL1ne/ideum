package users

import (
	"errors"
	"forum/internal/entity"
	"forum/internal/repository/users"
	"forum/internal/service/validator"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	SaveUser(*entity.UserSignupForm) (int, int, error)
	Authenticate(*entity.UserLoginForm) (int, int, error)
	GetAllUsers() ([]entity.UserEntity, error)
}

type userService struct {
	userRepo users.UserRepository
}

func NewUserService(u users.UserRepository) *userService {
	return &userService{
		userRepo: u,
	}
}

func (us *userService) SaveUser(u *entity.UserSignupForm) (int, int, error) {
	if !isRightSignUp(u) {
		return 0, http.StatusUnprocessableEntity, nil
	}

	id, err := us.userRepo.SaveUser(*u)
	if err != nil {
		if errors.Is(err, entity.ErrDuplicateEmail) {
			u.AddFieldError("email", "Email address is already in use")
			return 0, http.StatusUnprocessableEntity, err
		} else {
			return 0, http.StatusInternalServerError, err
		}
	}

	return id, http.StatusOK, nil
}

func (us *userService) Authenticate(u *entity.UserLoginForm) (int, int, error) {
	if !isRightLogin(u) {
		return 0, http.StatusUnprocessableEntity, nil
	}

	var userFromDB *entity.UserEntity
	var err error

	if validator.Matches(u.Identifier, EmailRX) {
		userFromDB, err = us.userRepo.GetUserByEmail(u.Identifier)
	} else {
		userFromDB, err = us.userRepo.GetUserByUsername(u.Identifier)
	}
	if err != nil {
		if errors.Is(err, entity.ErrInvalidCredentials) {
			return 0, http.StatusUnprocessableEntity, entity.ErrInvalidCredentials
		} else {
			return 0, http.StatusInternalServerError, err
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(userFromDB.Password), []byte(u.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, http.StatusUnprocessableEntity, entity.ErrInvalidCredentials
		} else {
			return 0, http.StatusInternalServerError, err
		}
	}

	return userFromDB.Id, http.StatusOK, nil
}

func (us *userService) GetAllUsers() ([]entity.UserEntity, error) {
	return nil, nil
}
