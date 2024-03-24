package entity

import "errors"

// Common errors
var (
	ErrNoRecord        = errors.New("entity: no matching row found")
	ErrInvalidFormData = errors.New("entity: some form data is invalid")
	ErrInvalidPathID   = errors.New("entity: invalid id in request path")
	ErrInvalidURLPath  = errors.New("entity: invalid url path")
)

// Post related errors
var (
	ErrInvalidPostID    = errors.New("entity: invalid post id")
	ErrInvalidImageType = errors.New("entity: invalid image type")
	ErrTooLargeImage    = errors.New("entity: too large image")
	ErrInvalidTags      = errors.New("entity: tags don't exist")
	ErrPostNotFound     = errors.New("entity: post not found")
	ErrCommentNotFound  = errors.New("entity: comment not found")
	ErrTagNotFound      = errors.New("entity: tag not found")
	ErrInvalidTag       = errors.New("entity: tag name is invalid")
	ErrDuplicateTag     = errors.New("entity: duplicate tag")
)

// User related errors
var (
	ErrDuplicateEmail        = errors.New("entity: duplicate email")
	ErrDuplicateUsername     = errors.New("entity: dupliate username")
	ErrInvalidCredentials    = errors.New("entity: invalid credentials")
	ErrInvalidUserID         = errors.New("entity: non-existent user id")
	ErrUnauthorized          = errors.New("entity: unauthorized")
	ErrForbiddenAccess       = errors.New("entity: further access is forbidden")
	ErrDuplicateNotification = errors.New("entity: duplicate notification")
	ErrNotificationNotFound  = errors.New("entity: notification not found")
	ErrDuplicatePromotion    = errors.New("entity: duplicate promotion")
	ErrPromotionNotFound     = errors.New("entity: promotion not found")
	ErrAdminNotFound         = errors.New("entity: admin not found")
	ErrDuplicateReport       = errors.New("entity: duplicate report")
	ErrReportNotFound        = errors.New("entity: report not found")
	ErrUserNotFound          = errors.New("entity: user not found")
)

// Notification related errors
var (
	ErrInvalidNotificaitonType = errors.New("entity: invalid notification type")
)
