package image

import (
	"database/sql"
)

type IImageRepository interface {
	SaveImage(int, string) error
	GetImage(int)(string,error)
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

func (r *imageRepository) SaveImage(postId int, imageName string) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	query := `
        INSERT INTO images (name, post_id)
        VALUES ($1, $2)
    `

	_, err = tx.Exec(query, imageName, postId)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r *imageRepository) GetImage(postID int) (string, error) {
	var imageName string
	query := `
        SELECT name FROM images WHERE post_id = $1
    `
	err := r.DB.QueryRow(query, postID).Scan(&imageName)
	switch {
	case err == sql.ErrNoRows: // image not found
		return "", nil
	case err != nil:
		return "", err
	}
	return imageName, nil
}
