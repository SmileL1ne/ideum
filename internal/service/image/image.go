package image

import (
	"forum/internal/repository/image"
)

type IImageService interface {
	SaveImage(int, string) error
	GetImage(int)(string,error)
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


func (is *imageService) SaveImage(postID int, name string) error {
	return is.imageRepo.SaveImage(postID, name)
}

func (is *imageService) GetImage(postID int)(string,error){
	return is.imageRepo.GetImage(postID)
}