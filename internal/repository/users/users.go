package users

import (
	"database/sql"
	"errors"
	"forum/internal/entity"

	"github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	SaveUser(entity.UserSignupForm) (int, error)
	GetUser(userId int) (entity.UserEntity, error)
	GetAllUsers() ([]entity.UserEntity, error)
}

type userRepository struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) *userRepository {
	return &userRepository{
		DB: db,
	}
}

func (r *userRepository) SaveUser(u entity.UserSignupForm) (int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	if err != nil {
		return 0, err
	}

	query := `INSERT INTO users (username, email, hashed_password, created) 
		VALUES ($1, $2, $3, datetime('now', 'utc', '+12 hours'))`

	result, err := r.DB.Exec(query, u.Username, u.Email, string(hashedPassword))
	if err != nil {
		var sqliteError *sqlite3.Error
		if errors.As(err, &sqliteError) {
			if sqliteError.Code == 19 {
				return 0, entity.ErrDuplicateEmail
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

func (r *userRepository) GetUser(userId int) (entity.UserEntity, error) {
	return entity.UserEntity{}, nil
}

func (r *userRepository) GetAllUsers() ([]entity.UserEntity, error) {
	return nil, nil
}
