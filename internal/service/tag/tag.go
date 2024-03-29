package tag

import (
	"forum/internal/entity"
	"forum/internal/repository/tag"
	"strconv"
)

type ITagService interface {
	GetAllTags() (*[]entity.TagEntity, error)
	AreTagsExist([]string) (bool, error)
	IsExist(int) (bool, error)
	DeleteTag(tagID int) error
	CreateTag(tag string) error
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
	return ts.tagRepo.GetAll()
}

func (ts *tagService) AreTagsExist(tags []string) (bool, error) {
	var tagIDs []int
	for _, tagIDStr := range tags {
		tagID, err := strconv.Atoi(tagIDStr)
		if err != nil {
			return false, entity.ErrInvalidTags
		}
		tagIDs = append(tagIDs, tagID)
	}

	return ts.tagRepo.AreTagsExist(tagIDs)
}

func (ts *tagService) IsExist(id int) (bool, error) {
	return ts.tagRepo.IsExist(id)
}

func (ts *tagService) DeleteTag(tagID int) error {
	return ts.tagRepo.Delete(tagID)
}

func (ts *tagService) CreateTag(tag string) error {
	if !IsValidTag(tag) {
		return entity.ErrInvalidTag
	}

	return ts.tagRepo.Create(tag)
}
