package user

import (
	"forum/internal/entity"
	"forum/internal/repository/user"
	"time"
)

var mockUser = entity.UserEntity{
	Id:        1,
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

func (r *UserRepoMock) SaveUser(u entity.UserSignupForm) (int, error) {
	if u.Username == "satoru" {
		return 0, entity.ErrDuplicateUsername
	} else if u.Email == "satoru@gmail.com" {
		return 0, entity.ErrDuplicateEmail
	}
	return mockUser.Id, nil
}

func (r *UserRepoMock) GetUserByEmail(email string) (entity.UserEntity, error) {
	if email == "yuta@gmail.com" {
		return mockUser, nil
	}
	return entity.UserEntity{}, entity.ErrInvalidCredentials
}

func (r *UserRepoMock) GetUserByUsername(username string) (entity.UserEntity, error) {
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
