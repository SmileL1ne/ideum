package tag

import (
	"database/sql"
	"forum/internal/entity"
	"sync"
)

type ITagRepository interface {
	GetAllTags() (*[]entity.TagEntity, error)
	AreTagsExist([]int) (bool, error)
	IsExist(int) (bool, error)
}

type tagRepo struct {
	DB *sql.DB
}

var _ ITagRepository = (*tagRepo)(nil)

func NewTagRepo(db *sql.DB) *tagRepo {
	return &tagRepo{
		DB: db,
	}
}

func (r *tagRepo) GetAllTags() (*[]entity.TagEntity, error) {
	query := `
		SELECT *
		FROM tags
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}

	var tags []entity.TagEntity
	for rows.Next() {
		var tag entity.TagEntity
		if err := rows.Scan(&tag.ID, &tag.Name, &tag.CreatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return &tags, nil
}

func (r *tagRepo) AreTagsExist(tagIDs []int) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT true
			FROM tags
			WHERE tags.id = $1
		)
	`

	var exists = true

	var wg sync.WaitGroup
	var errCh = make(chan error, len(tagIDs))

	wg.Add(len(tagIDs))
	for i, id := range tagIDs {
		go func(tagID int, it int) {
			defer wg.Done()

			var tagExists bool

			if err := r.DB.QueryRow(query, tagID).Scan(&tagExists); err != nil {
				errCh <- err
				return
			}
			exists = exists && tagExists
		}(id, i)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			return false, err
		}
	}

	return exists, nil
}

func (r *tagRepo) IsExist(id int) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT true
			FROM tags
			WHERE tags.id = $1
		)
	`

	var exists bool

	err := r.DB.QueryRow(query, id).Scan(&exists)

	return exists, err
}
