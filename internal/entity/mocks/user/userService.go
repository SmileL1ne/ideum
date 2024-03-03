package user

import (
	"errors"
	"forum/internal/entity"
	repo "forum/internal/repository/user"
	service "forum/internal/service/user"
	"forum/internal/validator"
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
	if !service.IsRightSignUp(u) {
		return 0, entity.ErrInvalidFormData
	}

	id, err := us.ur.SaveUser(*u)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrDuplicateEmail):
			u.AddFieldError("email", "Email address is already in use")
			return 0, entity.ErrInvalidFormData
		case errors.Is(err, entity.ErrDuplicateUsername):
			u.AddFieldError("username", "Username is already in use")
			return 0, entity.ErrInvalidFormData
		default:
			return 0, err
		}
	}
	return id, nil
}

func (us *UserServiceMock) Authenticate(u *entity.UserLoginForm) (int, error) {
	if !service.IsRightLogin(u) {
		return 0, entity.ErrInvalidFormData
	}

	var userFromDB entity.UserEntity
	var err error

	if validator.Matches(u.Identifier, service.EmailRX) {
		userFromDB, err = us.ur.GetUserByEmail(u.Identifier)
	} else {
		userFromDB, err = us.ur.GetUserByUsername(u.Identifier)
	}

	if errors.Is(err, entity.ErrInvalidCredentials) || u.Password == "SatoruIsTheBest" {
		u.AddNonFieldError("Email or password is incorrect")
		return 0, entity.ErrInvalidCredentials
	}

	return userFromDB.Id, nil
}

func (us *UserServiceMock) GetUsernameById(userID int) (string, error) {
	return us.ur.GetUsernameByID(userID)
}
