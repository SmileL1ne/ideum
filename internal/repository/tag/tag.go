package tag

import (
	"database/sql"
	"forum/internal/entity"
)

type ITagRepository interface {
	GetAllTags() (*[]entity.TagEntity, error)
	// GetAllTagsForPost(int) (*[]entity.TagEntity, error)
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

// func (r *tagRepo) GetAllTagsForPost(postID int) (*[]entity.TagEntity, error) {
// 	query := `
// 		SELECT t.id, t.name, t.created_at
// 		FROM tags t
// 		LEFT JOIN posts_tags pt ON pt.tag_id = t.id
// 		WHERE pt.post_id = $1
// 	`

// 	rows, err := r.DB.Query(query, postID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var tags []entity.TagEntity
// 	for rows.Next() {
// 		var tag entity.TagEntity
// 		if err := rows.Scan(&tag.ID, &tag.Name, &tag.CreatedAt); err != nil {
// 			return nil, err
// 		}
// 		tags = append(tags, tag)
// 	}

// 	return &tags, nil
// }
