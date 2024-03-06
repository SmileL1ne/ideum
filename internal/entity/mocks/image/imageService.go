package image

import (
	repo "forum/internal/repository/image"
	service "forum/internal/service/image"
	"mime/multipart"
)

type ImageServiceMock struct {
	ir repo.IImageRepository
}

func NewImageServiceMock(r repo.IImageRepository) *ImageServiceMock {
	return &ImageServiceMock{
		ir: r,
	}
}

var _ service.IImageService = (*ImageServiceMock)(nil)

// TODO: Add mock checks here (invalid image extension, ivalid image size)
func (is *ImageServiceMock) ProcessImage(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	return "mockImage", nil
}

func (is *ImageServiceMock) Get(postID int) (string, error) {
	return is.ir.ReadName(postID)
}
