package user

import (
	"errors"
	"forum/internal/entity"
	repo "forum/internal/repository/user"
	service "forum/internal/service/user"
)

type UserServiceMock struct {
	ur repo.IUserRepository
}

func NewUserServiceMock(r repo.IUserRepository) *UserServiceMock {
	return &UserServiceMock{
		ur: r,
	}
}

var _ service.IUserService = (*UserServiceMock)(nil)

func (us *UserServiceMock) SaveUser(u *entity.UserSignupForm) (int, error) {
	if u.Username == "satoru" && u.Email == "satoru@gmail.com" {
		return 0, entity.ErrInvalidFormData
	}

	id, err := us.ur.SaveUser(*u)
	if errors.Is(err, entity.ErrDuplicateEmail) {
		return 0, entity.ErrInvalidFormData
	} else if errors.Is(err, entity.ErrDuplicateUsername) {
		return 0, entity.ErrInvalidFormData
	}

	return id, nil
}

func (us *UserServiceMock) Authenticate(u *entity.UserLoginForm) (int, error) {
	if u.Identifier == "satoru" || u.Identifier == "satoru@gmail.com" {
		return 0, entity.ErrInvalidFormData
	}

	if u.Password == "satoruNumberOne" {
		return 0, entity.ErrInvalidCredentials
	}

	return mockUser.Id, nil // Avoid getting to repository, because it would return the same id, but more complicated way
}

func (us *UserServiceMock) GetUsernameById(userID int) (string, error) {
	return us.ur.GetUsernameByID(userID)
}
