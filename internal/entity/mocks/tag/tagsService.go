package tag

import (
	"forum/internal/entity"
	repo "forum/internal/repository/tag"
	service "forum/internal/service/tag"
	"strconv"
)

type TagServiceMock struct {
	tr repo.ITagRepository
}

func NewTagServiceMock(r repo.ITagRepository) *TagServiceMock {
	return &TagServiceMock{
		tr: r,
	}
}

var _ service.ITagService = (*TagServiceMock)(nil)

func (ts *TagServiceMock) GetAllTags() (*[]entity.TagEntity, error) {
	return ts.tr.GetAllTags()
}

func (ts *TagServiceMock) AreTagsExist(tags []string) (bool, error) {
	var tagIDs []int
	for _, tagIDStr := range tags {
		tagID, err := strconv.Atoi(tagIDStr)
		if err != nil {
			return false, entity.ErrInvalidFormData
		}
		tagIDs = append(tagIDs, tagID)
	}
	return ts.tr.AreTagsExist(tagIDs)
}

func (ts *TagServiceMock) IsExist(id int) (bool, error) {
	return ts.tr.IsExist(id)
}
