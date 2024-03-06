package post

import (
	"database/sql"
	"errors"
	"forum/internal/entity"
)

type IPostRepository interface {
	SavePost(entity.PostCreateForm, []int) (int, error)
	GetPost(int) (entity.PostEntity, error)
	GetAllPosts() (*[]entity.PostEntity, error)
	GetAllPostsByTagId(int) (*[]entity.PostEntity, error)
	GetAllPostsByUserID(int) (*[]entity.PostEntity, error)
	GetAllPostsByUserReaction(int) (*[]entity.PostEntity, error)
	ExistsPost(int) (bool, error)
}

type postRepository struct {
	DB *sql.DB
}

func NewPostRepo(db *sql.DB) *postRepository {
	return &postRepository{
		DB: db,
	}
}

var _ IPostRepository = (*postRepository)(nil)

func (r *postRepository) SavePost(p entity.PostCreateForm, tagIDs []int) (int, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	posts := `
		INSERT INTO posts (title, content, user_id, created_at) 
		VALUES ($1, $2, $3, datetime('now', 'localtime'))
		RETURNING id
	`
	var postID int
	err = tx.QueryRow(posts, p.Title, p.Content, p.UserID).Scan(&postID)
	if err != nil {
		return 0, err
	}

	posts_tags := `
		INSERT INTO posts_tags (post_id, tag_id, created_at)
		VALUES ($1, $2, datetime('now', 'localtime'))
	`
	for _, tagID := range tagIDs {
		_, err := tx.Exec(posts_tags, postID, tagID)
		if err != nil {
			return 0, err
		}
	}

	images := `
		INSERT INTO images (name, post_id)
		VALUES ($1, $2)
	`
	_, err = tx.Exec(images, p.ImageName, postID)
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return int(postID), nil
}

func (r *postRepository) GetPost(postID int) (entity.PostEntity, error) {
	query := `
		SELECT p.id, p.title, p.content, p.created_at, u.username,
			SUM(CASE WHEN pr.is_like = true THEN 1 ELSE 0 END) as likes_count,
			SUM(CASE WHEN pr.is_like = false THEN 1 ELSE 0 END) as dislikes_count,
			(
				SELECT COUNT(*)
				FROM comments c
				WHERE c.post_id = p.id
			),
			(
				SELECT GROUP_CONCAT(t.name, ', ')
				FROM tags t
				LEFT JOIN posts_tags pt ON pt.tag_id = t.id
				WHERE pt.post_id = p.id
			),
			(
				SELECT name
				FROM images
				WHERE post_id = p.id
			)
		FROM posts p
		INNER JOIN users u ON p.user_id = u.id
		LEFT JOIN post_reactions pr ON p.id = pr.post_id
		WHERE p.id=$1
		GROUP BY p.id
		`

	var post entity.PostEntity
	var tags sql.NullString
	var imageName sql.NullString
	if err := r.DB.QueryRow(query, postID).Scan(&post.ID, &post.Title, &post.Content,
		&post.CreatedAt, &post.Username, &post.Likes, &post.Dislikes, &post.CommentsLen, &tags, &imageName); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return entity.PostEntity{}, entity.ErrNoRecord
		}
		return entity.PostEntity{}, err
	}

	if tags.Valid {
		post.PostTags = tags.String
	}
	if imageName.Valid {
		post.ImageName = imageName.String
	}

	return post, nil
}

func (r *postRepository) GetAllPosts() (*[]entity.PostEntity, error) {
	query := `
		SELECT p.id, p.title, p.content, p.created_at, u.username, 
			SUM(CASE WHEN pr.is_like = true THEN 1 ELSE 0 END) as likes_count,
			SUM(CASE WHEN pr.is_like = false THEN 1 ELSE 0 END) as dislikes_count,
			(
				SELECT COUNT(*)
				FROM comments c
				WHERE c.post_id = p.id
			),
			(
				SELECT GROUP_CONCAT(t.name, ', ')
				FROM tags t
				LEFT JOIN posts_tags pt ON pt.tag_id = t.id
				WHERE pt.post_id = p.id
			),
			(
				SELECT name
				FROM images
				WHERE post_id = p.id
			)
		FROM posts p
		INNER JOIN users u ON p.user_id = u.id
		LEFT JOIN post_reactions pr ON p.id = pr.post_id
		GROUP BY p.id
		ORDER BY p.created_at DESC
	`

	return getAllPostsByQuery(r.DB, query)
}

func (r *postRepository) GetAllPostsByTagId(tagID int) (*[]entity.PostEntity, error) {
	query := `
		SELECT p.id, p.title, p.content, p.created_at, u.username, 
			SUM(CASE WHEN pr.is_like = true THEN 1 ELSE 0 END) as likes_count,
			SUM(CASE WHEN pr.is_like = false THEN 1 ELSE 0 END) as dislikes_count,
			(
				SELECT COUNT(*)
				FROM comments c
				WHERE c.post_id = p.id
			),
			(
				SELECT GROUP_CONCAT(t.name, ', ')
				FROM tags t
				LEFT JOIN posts_tags pt ON pt.tag_id = t.id
				WHERE pt.post_id = p.id
			),
			(
				SELECT name
				FROM images
				WHERE post_id = p.id
			)
		FROM posts p
		INNER JOIN users u ON p.user_id = u.id
		LEFT JOIN post_reactions pr ON p.id = pr.post_id
		WHERE p.id IN (
			SELECT pt.post_id
			FROM posts_tags pt
			WHERE pt.tag_id = $1
			)
		GROUP BY p.id
		ORDER BY p.created_at DESC
	`

	return getAllPostsByQuery(r.DB, query, tagID)
}

func (r *postRepository) GetAllPostsByUserID(userID int) (*[]entity.PostEntity, error) {
	query := `
		SELECT p.id, p.title, p.content, p.created_at, u.username,
			SUM(CASE WHEN pr.is_like = true THEN 1 ELSE 0 END) as likes_count,
			SUM(CASE WHEN pr.is_like = false THEN 1 ELSE 0 END) as dislikes_count,
			(
				SELECT COUNT(*)
				FROM comments c
				WHERE c.post_id = p.id
			),
			(
				SELECT GROUP_CONCAT(t.name, ', ')
				FROM tags t
				LEFT JOIN posts_tags pt ON pt.tag_id = t.id
				WHERE pt.post_id = p.id
			),
			(
				SELECT name
				FROM images
				WHERE post_id = p.id
			)
		FROM posts p
		INNER JOIN users u ON p.user_id = u.id
		LEFT JOIN post_reactions pr ON p.id = pr.post_id
		WHERE p.user_id = $1
		GROUP BY p.id 
		ORDER BY p.created_at DESC
	`

	return getAllPostsByQuery(r.DB, query, userID)
}

func (r *postRepository) GetAllPostsByUserReaction(userID int) (*[]entity.PostEntity, error) {
	query := `
		SELECT p.id, p.title, p.content, p.created_at, u.username,
			SUM(CASE WHEN pr.is_like = true THEN 1 ELSE 0 END) as likes_count,
			SUM(CASE WHEN pr.is_like = false THEN 1 ELSE 0 END) as dislikes_count,
			(
				SELECT COUNT(*)
				FROM comments c
				WHERE c.post_id = p.id
			),
			(
				SELECT GROUP_CONCAT(t.name, ', ')
				FROM tags t
				LEFT JOIN posts_tags pt ON pt.tag_id = t.id
				WHERE pt.post_id = p.id 
			),
			(
				SELECT name
				FROM images
				WHERE post_id = p.id
			)
		FROM posts p
		INNER JOIN users u ON p.user_id = u.id
		LEFT JOIN post_reactions pr ON p.id = pr.post_id
		WHERE pr.user_id = $1
		GROUP BY p.id 
	`

	return getAllPostsByQuery(r.DB, query, userID)
}

func (r *postRepository) ExistsPost(postID int) (bool, error) {
	var exists bool

	query := `
		SELECT EXISTS(
			SELECT true
			FROM posts
			WHERE id = $1
		)
	`

	err := r.DB.QueryRow(query, postID).Scan(&exists)
	return exists, err
}
