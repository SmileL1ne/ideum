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
	GetReports() (*[]entity.Report, error)
	Promote(userID int) error
	Demote(userID int) error
	GetUsers() (*[]entity.UserEntity, error)
	FindNotification(nType string, userFrom, userTo int) (int, error)
	GetNotificationsCount(userID int) (int, error)
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

	err := r.DB.QueryRow(query, value).Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.CreatedAt)
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

	return err
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
			if sqliteError.Code == 19 && strings.Contains(sqliteError.Error(), "UNIQUE constraint failed:") {
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
			if sqliteError.Code == 19 && strings.Contains(sqliteError.Error(), "UNIQUE constraint failed:") {
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
		ORDER BY n.created_at DESC
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

func (r *userRepository) GetNotificationsCount(userID int) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM notifications 
		WHERE user_to = $1
	`

	var count int

	err := r.DB.QueryRow(query, userID).Scan(&count)

	return count, err
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
		SELECT r.id, r.user_id, r.created_at, u.username
		FROM requests r
		INNER JOIN users u ON r.user_id = u.id
		ORDER BY r.created_at DESC
	`

	var requests []entity.Request

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var req entity.Request
		if err := rows.Scan(&req.ID, &req.UserID, &req.CreatedAt, &req.Username); err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}

	return &requests, nil
}

func (r *userRepository) GetReports() (*[]entity.Report, error) {
	query := `
		SELECT r.id, r.reason, r.user_from, r.source_id, r.source_type, r.created_at, u.username
		FROM reports r
		INNER JOIN users u ON u.id = r.user_from
		ORDER BY r.created_at DESC
	`

	var reports []entity.Report

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r entity.Report
		if err := rows.Scan(&r.ID, &r.Reason, &r.UserFrom, &r.SourceID, &r.SourceType, &r.CreatedAt, &r.Username); err != nil {
			return nil, err
		}
		reports = append(reports, r)
	}

	return &reports, nil
}

func (r *userRepository) Promote(userID int) error {
	query := `
		UPDATE roles
		SET role = $1
		WHERE user_id = $2
	`

	res, err := r.DB.Exec(query, entity.MODERATOR, userID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return entity.ErrUserNotFound
	}

	return err
}

func (r *userRepository) Demote(userID int) error {
	query := `
		UPDATE roles
		SET role = $1
		WHERE user_id = $2
	`

	res, err := r.DB.Exec(query, entity.USER, userID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return entity.ErrUserNotFound
	}

	return err
}

func (r *userRepository) GetUsers() (*[]entity.UserEntity, error) {
	query := `
		SELECT u.id, u.username, u.email, u.hashed_password, u.created_at, r.role
		FROM users u
		LEFT JOIN roles r ON u.id = r.user_id
		ORDER BY u.created_at DESC
	`

	var users []entity.UserEntity

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u entity.UserEntity
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.CreatedAt, &u.Role); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return &users, nil
}

func (r *userRepository) FindNotification(nType string, userFrom, userTo int) (int, error) {
	query := `
		SELECT id
		FROM notifications
		WHERE type = $1 AND user_from = $2 AND user_to = $3
	`

	var notificationID int

	err := r.DB.QueryRow(query, nType, userFrom, userTo).Scan(&notificationID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, entity.ErrNotificationNotFound
		}
		return 0, err
	}

	return notificationID, nil
}
