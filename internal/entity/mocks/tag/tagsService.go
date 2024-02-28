package tag

import (
	"forum/internal/entity"
	repo "forum/internal/repository/tag"
	service "forum/internal/service/tag"
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
