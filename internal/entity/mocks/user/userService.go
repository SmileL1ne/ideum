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

	id, err := us.ur.Insert(*u, nil)
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
		userFromDB, err = us.ur.GetByEmail(u.Identifier)
	} else {
		userFromDB, err = us.ur.GetByUsername(u.Identifier)
	}

	if errors.Is(err, entity.ErrInvalidCredentials) || u.Password == "SatoruIsTheBest" {
		u.AddNonFieldError("Email or password is incorrect")
		return 0, entity.ErrInvalidCredentials
	}

	return userFromDB.ID, nil
}

func (us *UserServiceMock) GetUsernameById(userID int) (string, error) {
	return us.ur.GetUsernameByID(userID)
}

func (us *UserServiceMock) GetUserByEmail(email string) (entity.UserEntity, error) {
	return us.ur.GetByEmail(email)
}

func (us *UserServiceMock) GetUserRole(userID int) (string, error) {
	return us.ur.GetRole(userID)
}

func (us *UserServiceMock) SendNotification(notification entity.Notification) error {
	return us.ur.CreateNotification(notification)
}

func (us *UserServiceMock) GetRequests(role string) (*[]entity.Notification, error) {
	if role != entity.ADMIN {
		return nil, entity.ErrForbiddenAccess
	}

	return us.ur.GetRequests()
}

func (us *UserServiceMock) PromoteUser(userID int) error {
	return us.ur.Promote(userID)
}
