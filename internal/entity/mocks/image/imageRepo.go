package image

import (
	"forum/internal/entity"
	"forum/internal/repository/image"
)

type ImageRepoMock struct {
}

func NewImageRepoMock() *ImageRepoMock {
	return &ImageRepoMock{}
}

var _ image.IImageRepository = (*ImageRepoMock)(nil)

func (r *ImageRepoMock) Create(i entity.ImageEntity) error {
	return nil
}

func (r *ImageRepoMock) ReadName(postID int) (string, error) {
	if postID == 0 {
		return "mockImage", nil
	}
	return "", nil
}
