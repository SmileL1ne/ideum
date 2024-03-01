package user

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/entity"
	"strings"

	"github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type IUserRepository interface {
	SaveUser(entity.UserSignupForm) (int, error)
	GetUserByUsername(username string) (entity.UserEntity, error)
	GetUserByEmail(email string) (entity.UserEntity, error)
	GetUsernameByID(userID int) (string, error)
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

func (r *userRepository) SaveUser(u entity.UserSignupForm) (int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	if err != nil {
		return 0, err
	}

	user, err := r.GetUserByUsername(u.Username)
	if err != nil && !errors.Is(err, entity.ErrInvalidCredentials) {
		return 0, err
	}
	if user != (entity.UserEntity{}) {
		if strings.EqualFold(u.Username, user.Username) {
			return 0, entity.ErrDuplicateUsername
		}
	}

	query := `INSERT INTO users (username, email, hashed_password, created_at) 
		VALUES ($1, $2, $3, datetime('now', 'localtime'))`

	result, err := r.DB.Exec(query, u.Username, u.Email, string(hashedPassword))
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

	return int(id), nil
}

func (r *userRepository) GetUserByEmail(email string) (entity.UserEntity, error) {
	return r.getUserByField("email", email)
}

func (r *userRepository) GetUserByUsername(username string) (entity.UserEntity, error) {
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
