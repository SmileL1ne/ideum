package user

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/entity"
	"strings"

	"github.com/mattn/go-sqlite3"
)

type IUserRepository interface {
	Insert(user entity.UserSignupForm, hashedPassword []byte) (int, error)
	GetByUsername(username string) (entity.UserEntity, error)
	GetByEmail(email string) (entity.UserEntity, error)
	GetUsernameByID(userID int) (string, error)
	GetRole(userID int) (string, error)
	CreateNotification(n entity.Notification) error
	CreatePromotion(userID int) error
	CreateReport(report entity.Report) error
	DeleteReport(reportID int) error
	DeletePromotion(promotionID int) error
	GetNotifications(userID int) (*[]entity.Notification, error)
	DeleteNotification(notificationID int) error
	GetRequests() (*[]entity.Request, error)
	Promote(userID int) error
}

type userRepository struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) *userRepository {
	return &userRepository{
		DB: db,
	}
}

var _ IUserRepository = (*userRepository)(nil)

func (r *userRepository) Insert(u entity.UserSignupForm, hashedPassword []byte) (int, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	query1 := `INSERT INTO users (username, email, hashed_password, created_at) 
		VALUES ($1, $2, $3, datetime('now', 'localtime'))`

	result, err := tx.Exec(query1, u.Username, u.Email, string(hashedPassword))
	if err != nil {
		var sqliteError sqlite3.Error
		if errors.As(err, &sqliteError) {
			if sqliteError.Code == 19 && strings.Contains(sqliteError.Error(), "UNIQUE constraint failed:") {
				switch {
				case strings.Contains(sqliteError.Error(), "users.email"):
					return 0, entity.ErrDuplicateEmail
				case strings.Contains(sqliteError.Error(), "users.username"):
					return 0, entity.ErrDuplicateUsername
				default:
					return 0, fmt.Errorf("(repo) SaveUser: unknown field - %v", sqliteError)
				}
			}
		}
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	query2 := `
		INSERT INTO roles (role, user_id)
		VALUES ($1, $2)
	`

	_, err = tx.Exec(query2, entity.USER, int(id))
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *userRepository) GetByEmail(email string) (entity.UserEntity, error) {
	return r.getUserByField("email", email)
}

func (r *userRepository) GetByUsername(username string) (entity.UserEntity, error) {
	return r.getUserByField("username", username)
}

func (r *userRepository) GetUsernameByID(userID int) (string, error) {
	user, err := r.getUserByField("id", userID)
	return user.Username, err
}

func (r *userRepository) getUserByField(field string, value interface{}) (entity.UserEntity, error) {
	var u entity.UserEntity

	query := fmt.Sprintf(`SELECT * FROM users WHERE %s = $1 COLLATE NOCASE`, field)

	err := r.DB.QueryRow(query, value).Scan(&u.Id, &u.Username, &u.Email, &u.Password, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.UserEntity{}, entity.ErrInvalidCredentials
		}
		return entity.UserEntity{}, err
	}

	return u, nil
}

func (r *userRepository) GetRole(userID int) (string, error) {
	query := `
		SELECT role
		FROM roles
		WHERE user_id = $1
	`

	var role string

	err := r.DB.QueryRow(query, userID).Scan(&role)
	if err != nil {
		return "", err
	}

	return role, nil
}

func (r *userRepository) CreateNotification(n entity.Notification) error {
	query := `
		INSERT INTO notifications (type, user_from, user_to, content, source_id, source_type, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, datetime('now', 'localtime'))
	`

	_, err := r.DB.Exec(query, n.Type, n.UserFrom, n.UserTo, n.Content, n.SourceID, n.SourceType)
	if err != nil {
		var sqliteError sqlite3.Error
		if errors.As(err, &sqliteError) {
			if sqliteError.Code == 19 && strings.Contains(sqliteError.Error(), "UNIQUE constraint failed:") {
				switch {
				case strings.Contains(sqliteError.Error(), "notifications.type, notifications.source_id, notifications.user_from, notifications.user_to"):
					return entity.ErrDuplicateNotification
				default:
					return fmt.Errorf("(repo) SaveUser: unknown field - %v", sqliteError)
				}
			}
		}
		return err
	}

	return nil
}

func (r *userRepository) CreatePromotion(userID int) error {
	query := `
		INSERT INTO requests (user_id, created_at)
		VALUES ($1, datetime('now', 'localtime'))
	`

	_, err := r.DB.Exec(query, userID)
	if err != nil {
		var sqliteError sqlite3.Error
		if errors.As(err, &sqliteError) {
			if sqliteError.Code == 19 && strings.Contains(sqliteError.Error(), "UNIQUE constaint failed:") {
				return entity.ErrDuplicatePromotion
			}
		}
		return err
	}

	return nil
}

func (r *userRepository) CreateReport(report entity.Report) error {
	query := `
		INSERT INTO reports (reason, user_from, source_id, source_type, created_at)
		VALUES ($1, $2, $3, $4, datetime('now', 'localtime'))
	`

	_, err := r.DB.Exec(query, report.Reason, report.UserFrom, report.SourceID, report.SourceType)
	if err != nil {
		var sqliteError sqlite3.Error
		if errors.As(err, &sqliteError) {
			if sqliteError.Code == 19 && strings.Contains(sqliteError.Error(), "UNIQUE constaint failed:") {
				return entity.ErrDuplicateReport
			}
		}
		return err
	}

	return nil
}

func (r *userRepository) DeleteReport(reportID int) error {
	query := `
		DELETE FROM reports
		WHERE id = $1
	`

	res, err := r.DB.Exec(query, reportID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return entity.ErrReportNotFound
	}

	return nil
}

func (r *userRepository) DeletePromotion(promotionID int) error {
	query := `
		DELETE FROM requests
		WHERE id = $1
	`

	res, err := r.DB.Exec(query, promotionID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return entity.ErrPromotionNotFound
	}

	return nil
}

func (r *userRepository) GetNotifications(userID int) (*[]entity.Notification, error) {
	query := `
		SELECT n.ID, n.type, n.user_from, n.user_to, n.content, n.source_id, n.source_type, n.created_at, u.username
		FROM notifications n
		JOIN users u ON u.id = n.user_from
		WHERE user_to = $1
	`

	var notifications []entity.Notification

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var n entity.Notification
		if err := rows.Scan(&n.ID, &n.Type, &n.UserFrom, &n.UserTo, &n.Content, &n.SourceID, &n.SourceType, &n.CreatedAt, &n.Username); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}

	return &notifications, nil
}

func (r *userRepository) DeleteNotification(notificationID int) error {
	query := `
		DELETE FROM notifications
		WHERE id = $1
	`

	res, err := r.DB.Exec(query, notificationID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return entity.ErrNotificationNotFound
	}

	return nil
}

func (r *userRepository) GetRequests() (*[]entity.Request, error) {
	query := `
		SELECT r.id, r.user_id, u.username
		FROM requests r
		INNER JOIN users u ON r.user_id = u.id
	`

	var requests []entity.Request

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var req entity.Request
		if err := rows.Scan(&req.ID, &req.UserID, &req.Username, &req.CreatedAt); err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}

	return &requests, nil
}

func (r *userRepository) Promote(userID int) error {
	query := `
		UPDATE roles
		SET role = $1
		WHERE user_id = $2
	`

	_, err := r.DB.Exec(query, entity.MODERATOR, userID)

	return err
}
