package user

import (
	"errors"
	"forum/internal/entity"
	"forum/internal/repository/user"
	"forum/internal/validator"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	SaveUser(*entity.UserSignupForm) (int, error)
	Authenticate(*entity.UserLoginForm) (int, error)
	GetUsernameById(int) (string, error)
	GetUserByEmail(string) (entity.UserEntity, error)
	GetUserRole(int) (string, error)
	SendNotification(notification entity.Notification) error
	SendPromotion(userID int) error
	SendReport(report entity.Report) error
	DeleteReport(reportID int) error
	DeletePromotion(promotionID int) error
	GetRequests() (*[]entity.Request, error)
	GetReports() (*[]entity.Report, error)
	PromoteUser(userID int) error
	GetNotifications(userID int) (*[]entity.Notification, error)
	DeleteNotification(notificationID int) error
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
	case entity.PROMOTED:
		n.Content = "Congratulations! You've been promoted to moderator!"
	case entity.POST_LIKE:
		n.Content = "Liked your post"
	case entity.POST_DISLIKE:
		n.Content = "Disliked your post"
	case entity.COMMENT_LIKE:
		n.Content = "Liked your comment"
	case entity.COMMENT_DISLIKE:
		n.Content = "Disliked your comment"
	case entity.COMMENTED:
		n.Content = "Left a comment on your post"
	case entity.REJECT_PROMOTION:
		n.Content = "Your promotion was rejected"
	case entity.REJECT_REPORT:
		n.Content = "Your report was rejected"
	case entity.DELETE_POST:
		n.Content = "Your post/posts was/were deleted" + n.Content
	case entity.DELETE_COMMENT:
		n.Content = "Your comment/comments was/were deleted" + n.Content
	default:
		return entity.ErrInvalidNotificaitonType
	}

	return us.userRepo.CreateNotification(n)
}

func (us *userService) SendPromotion(userID int) error {
	return us.userRepo.CreatePromotion(userID)
}

func (us *userService) SendReport(report entity.Report) error {
	return us.userRepo.CreateReport(report)
}

func (us *userService) DeleteReport(reportID int) error {
	return us.userRepo.DeleteReport(reportID)
}

func (us *userService) DeletePromotion(promotionID int) error {
	return us.userRepo.DeletePromotion(promotionID)
}

func (us *userService) GetRequests() (*[]entity.Request, error) {
	return us.userRepo.GetRequests()
}

func (us *userService) GetReports() (*[]entity.Report, error) {
	return us.userRepo.GetReports()
}

func (us *userService) GetNotifications(userID int) (*[]entity.Notification, error) {
	return us.userRepo.GetNotifications(userID)
}

func (us *userService) PromoteUser(userID int) error {
	return us.userRepo.Promote(userID)
}

func (us *userService) DeleteNotification(notificationID int) error {
	return us.userRepo.DeleteNotification(notificationID)
}
