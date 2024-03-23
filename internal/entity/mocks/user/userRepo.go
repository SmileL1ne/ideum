package user

import (
	"forum/internal/entity"
	"forum/internal/repository/user"
	"time"
)

var mockUser = entity.UserEntity{
	ID:        1,
	Username:  "yuta",
	Email:     "yuta@gmail.com",
	Password:  "yuta12345",
	CreatedAt: time.Date(2003, time.July, 6, 0, 0, 0, 0, time.Local),
}

type UserRepoMock struct{}

func NewUserRepoMock() *UserRepoMock {
	return &UserRepoMock{}
}

var _ user.IUserRepository = (*UserRepoMock)(nil)

func (r *UserRepoMock) Insert(u entity.UserSignupForm, hashedPassword []byte) (int, error) {
	if u.Username == "satoru" {
		return 0, entity.ErrDuplicateUsername
	} else if u.Email == "satoru@gmail.com" {
		return 0, entity.ErrDuplicateEmail
	}
	return mockUser.ID, nil
}

func (r *UserRepoMock) GetByEmail(email string) (entity.UserEntity, error) {
	if email == "yuta@gmail.com" {
		return mockUser, nil
	}
	return entity.UserEntity{}, entity.ErrInvalidCredentials
}

func (r *UserRepoMock) GetByUsername(username string) (entity.UserEntity, error) {
	if username == "yuta" {
		return mockUser, nil
	}
	return entity.UserEntity{}, entity.ErrInvalidCredentials
}

func (r *UserRepoMock) GetUsernameByID(userID int) (string, error) {
	if userID == 2 {
		return "", entity.ErrInvalidCredentials
	}
	return mockUser.Username, nil
}

func (r *UserRepoMock) GetRole(userID int) (string, error) {
	switch userID {
	case 1:
		return "admin", nil
	case 2:
		return "moderator", nil
	default:
		return "guest", nil
	}
}

func (r *UserRepoMock) CreateNotification(n entity.Notification) error {
	return nil
}

func (r *UserRepoMock) GetRequests() (*[]entity.Notification, error) {
	return nil, nil
}

func (r *UserRepoMock) Promote(userID int) error {
	return nil
}
