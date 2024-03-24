package image

import (
	"database/sql"
	"forum/internal/entity"
)

type IImageRepository interface {
	Create(entity.ImageEntity) error
	GetName(int) (string, error)
}

type imageRepository struct {
	DB *sql.DB
}

var _ IImageRepository = (*imageRepository)(nil)

func NewImageRepo(db *sql.DB) *imageRepository {
	return &imageRepository{
		DB: db,
	}
}

func (r *imageRepository) Create(i entity.ImageEntity) error {
	query := `
		INSERT INTO images (name, post_id)
		VALUES ($1, $2)
	`

	_, err := r.DB.Exec(query, i.Name, i.PostID)

	return err
}

func (r *imageRepository) GetName(postID int) (string, error) {
	query := `
		SELECT (name)
		FROM images
		WHERE post_id = $1
	`

	var name string

	err := r.DB.QueryRow(query, postID).Scan(&name)

	return name, err
}
