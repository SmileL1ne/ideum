package users

import (
	"errors"
	"forum/internal/entity"
	"forum/internal/repository/users"
	"net/http"
)

type UserService interface {
	SaveUser(*entity.UserSignupForm) (int, int, error)
	GetUser(userId string) (entity.UserEntity, int, error)
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
	if !isRightUser(u) {
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

func (us *userService) GetUser(userId string) (entity.UserEntity, int, error) {
	return entity.UserEntity{}, 0, nil
}

func (us *userService) GetAllUsers() ([]entity.UserEntity, error) {
	return nil, nil
}
