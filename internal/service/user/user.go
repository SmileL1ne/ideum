package user

import (
	"errors"
	"forum/internal/entity"
	"forum/internal/repository/user"
	"forum/internal/validator"
	"strings"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	SaveUser(*entity.UserSignupForm) (int, error)
	Authenticate(*entity.UserLoginForm) (int, error)
	GetUsernameById(int) (string, error)
	GetUserByEmail(string) (entity.UserEntity, error)
	GetUserRole(int) (string, error)
	SendNotification(notification entity.Notification) error
	GetRequests(role string) (*[]entity.Notification, error)
	PromoteUser(userID int) error
}

type userService struct {
	userRepo user.IUserRepository
}

func NewUserService(u user.IUserRepository) *userService {
	return &userService{
		userRepo: u,
	}
}

var _ IUserService = (*userService)(nil)

func (us *userService) SaveUser(u *entity.UserSignupForm) (int, error) {
	if !IsRightSignUp(u) {
		return 0, entity.ErrInvalidFormData
	}

	// Check if used with that username doesn't already exist
	user, err := us.userRepo.GetByUsername(u.Username)
	if err != nil && !errors.Is(err, entity.ErrInvalidCredentials) {
		return 0, err
	}
	if user != (entity.UserEntity{}) && strings.EqualFold(u.Username, user.Username) {
		return 0, entity.ErrDuplicateUsername
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	if err != nil {
		return 0, err
	}

	id, err := us.userRepo.Insert(*u, hashedPassword)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrDuplicateEmail):
			u.AddFieldError("email", "Email address is already in use")
			return 0, entity.ErrDuplicateEmail
		case errors.Is(err, entity.ErrDuplicateUsername):
			u.AddFieldError("username", "Username is already in use")
			return 0, entity.ErrDuplicateUsername
		default:
			return 0, err
		}
	}

	return id, nil
}

func (us *userService) Authenticate(u *entity.UserLoginForm) (int, error) {
	if !IsRightLogin(u) {
		return 0, entity.ErrInvalidFormData
	}

	var userFromDB entity.UserEntity
	var err error

	if validator.Matches(u.Identifier, EmailRX) {
		userFromDB, err = us.userRepo.GetByEmail(u.Identifier)
	} else {
		userFromDB, err = us.userRepo.GetByUsername(u.Identifier)
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

func (us *userService) GetUsernameById(userID int) (string, error) {
	return us.userRepo.GetUsernameByID(userID)
}

func (us *userService) GetUserByEmail(email string) (entity.UserEntity, error) {
	return us.userRepo.GetByEmail(email)
}

func (us *userService) GetUserRole(userID int) (string, error) {
	return us.userRepo.GetRole(userID)
}

func (us *userService) SendNotification(n entity.Notification) error {
	switch n.Type {
	case entity.PROMOTION:
		n.Content = "requested promotion to moderator"
	case entity.POST_LIKE:
		n.Content = "liked your post"
	case entity.POST_DISLIKE:
		n.Content = "disliked your post"
	case entity.COMMENT_LIKE:
		n.Content = "liked your comment"
	case entity.COMMENT_DISLIKE:
		n.Content = "disliked your comment"
	case entity.COMMENTED:
		n.Content = "left a comment on your post"
	case entity.REPORT:
		n.Content = "reported this content as " + n.Content
	case entity.REJECT_PROMOTION:
		n.Content = "rejected your promotion to moderator"
	default:
		return entity.ErrInvalidNotificaitonType
	}

	return us.userRepo.CreateNotification(n)
}

func (us *userService) GetRequests(role string) (*[]entity.Notification, error) {
	if role != entity.ADMIN {
		return nil, entity.ErrForbiddenAccess
	}

	requests, err := us.userRepo.GetRequests()
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	var errCh = make(chan error, len(*requests))

	for _, req := range *requests {
		wg.Add(1)
		go func(request entity.Notification) {
			username, err := us.userRepo.GetUsernameByID(request.UserFrom)
			if err != nil {
				errCh <- err
				return
			}
			request.Username = username
		}(req)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		return nil, err
	}

	return requests, nil
}

func (us *userService) PromoteUser(userID int) error {
	return us.userRepo.Promote(userID)
}
