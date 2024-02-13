package tag

import (
	"forum/internal/entity"
	"forum/internal/repository/tag"
)

type ITagService interface {
	GetAllTags() (*[]entity.TagEntity, error)
	GetAllTagsForPost(int) (*[]entity.TagEntity, error)
}

type tagService struct {
	tagRepo tag.ITagRepository
}

var _ ITagService = (*tagService)(nil)

func NewTagService(r tag.ITagRepository) *tagService {
	return &tagService{
		tagRepo: r,
	}
}

func (ts *tagService) GetAllTags() (*[]entity.TagEntity, error) {
	return ts.tagRepo.GetAllTags()
}

func (ts *tagService) GetAllTagsForPost(postID int) (*[]entity.TagEntity, error) {
	return ts.tagRepo.GetAllTagsForPost(postID)
}