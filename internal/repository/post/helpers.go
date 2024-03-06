package post

import (
	"database/sql"
	"forum/internal/entity"
)

func getAllPostsByQuery(db *sql.DB, query string, args ...interface{}) (*[]entity.PostEntity, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	var posts []entity.PostEntity

	for rows.Next() {
		var post entity.PostEntity
		var tags sql.NullString
		var imageName sql.NullString
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt,
			&post.Username, &post.Likes, &post.Dislikes, &post.CommentsLen, &tags, &imageName); err != nil {

			return nil, err
		}
		if tags.Valid {
			post.PostTags = tags.String
		}
		if imageName.Valid {
			post.ImageName = imageName.String
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &posts, nil
}
