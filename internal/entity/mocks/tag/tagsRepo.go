package tag

import (
	"forum/internal/entity"
	"forum/internal/repository/tag"
	"time"
)

var mockTag = entity.TagEntity{
	ID:        1,
	Name:      "Art",
	CreatedAt: time.Date(2003, time.July, 6, 0, 0, 0, 0, time.Local),
}

type TagRepoMock struct{}

func NewTagRepoMock() *TagRepoMock {
	return &TagRepoMock{}
}

var _ tag.ITagRepository = (*TagRepoMock)(nil)

func (r *TagRepoMock) GetAllTags() (*[]entity.TagEntity, error) {
	return &[]entity.TagEntity{mockTag}, nil
}

func (r *TagRepoMock) AreTagsExist(tagIDs []int) (bool, error) {
	if tagIDs[0] == 1 {
		return true, nil
	}
	return false, nil
}

func (r *TagRepoMock) IsExists(id int) (bool, error) {
	if id == 1 {
		return true, nil
	}
	return false, nil
}
