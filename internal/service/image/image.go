package image

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/repository/image"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

type IImageService interface {
	ProcessImage(multipart.File, *multipart.FileHeader) (string, error)
	Get(int) (string, error)
}

type imageService struct {
	imageRepo image.IImageRepository
}

var _ IImageService = (*imageService)(nil)

func NewImageService(r image.IImageRepository) *imageService {
	return &imageService{
		imageRepo: r,
	}
}

// ProcessImage processes given image file.
//
// It returns new name of the image file created by hashing the contents of the file
// (to prevent saving duplicate images)
func (is *imageService) ProcessImage(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	defer file.Close()

	ext := strings.Split(fileHeader.Filename, ".")[1]
	h := sha256.New()
	_, err := io.Copy(h, file)
	if err != nil {
		return "", err
	}

	fname := fmt.Sprintf("%x.%s", h.Sum(nil), ext)
	path := filepath.Join("./web/static/public/", fname)
	newFile, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer newFile.Close()

	file.Seek(0, 0)
	_, err = io.Copy(newFile, file)
	return fname, err
}

func (is *imageService) Get(postID int) (string, error) {
	name, err := is.imageRepo.ReadName(postID)
	if errors.Is(err, sql.ErrNoRows) {
		return name, nil
	}
	return name, err
}
