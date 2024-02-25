package image

import (
	"forum/internal/repository/image"
)

type IImageService interface {
	SaveImage(int, string) error
}

type imageService struct {
	imageRepo image.IImageRepository
}

var _ IImageService = (*imageService)(nil)

func (is *imageService) SaveImage(postID int, name string) error {
	return is.imageRepo.SaveImage(postID, name)
}

func NewImageService(r image.IImageRepository) *imageService {
	return &imageService{
		imageRepo: r,
	}
}
