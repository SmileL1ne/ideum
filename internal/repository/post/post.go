package post

import (
	"database/sql"
	"errors"
	"forum/internal/entity"
)

type IPostRepository interface {
	Insert(entity.PostCreateForm, []int) (int, error)
	Get(int) (entity.PostEntity, error)
	GetAll() (*[]entity.PostEntity, error)
	GetAllByTagId(int) (*[]entity.PostEntity, error)
	GetAllByUserID(int) (*[]entity.PostEntity, error)
	GetAllByUserReaction(int) (*[]entity.PostEntity, error)
	GetAllCommentedPosts(userID int) (*[]entity.PostEntity, error)
	Exists(int) (bool, error)
	Delete(postID int, userID int) error
	DeleteByPrivileged(postID int) error
	GetAuthorID(postID int) (int, error)
	Update(p entity.PostCreateForm, tagIDs []int, deleteImage bool) error
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

func (r *postRepository) Insert(p entity.PostCreateForm, tagIDs []int) (int, error) {
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

func (r *postRepository) Get(postID int) (entity.PostEntity, error) {
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

func (r *postRepository) GetAll() (*[]entity.PostEntity, error) {
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

func (r *postRepository) GetAllByTagId(tagID int) (*[]entity.PostEntity, error) {
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

func (r *postRepository) GetAllByUserID(userID int) (*[]entity.PostEntity, error) {
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

func (r *postRepository) GetAllByUserReaction(userID int) (*[]entity.PostEntity, error) {
	query := `
		SELECT p.id, p.title, p.content, p.created_at, u.username,
			COALESCE(l.likes_count, 0) as likes_count,
			COALESCE(d.dislikes_count, 0) as dislikes_count,
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
		LEFT JOIN (
			SELECT post_id, COUNT(*) as likes_count
			FROM post_reactions
			WHERE is_like = true
			GROUP BY post_id
		) l ON p.id = l.post_id
		LEFT JOIN (
			SELECT post_id, COUNT(*) as dislikes_count
			FROM post_reactions
			WHERE is_like = false
			GROUP BY post_id
		) d ON p.id = d.post_id
		WHERE p.id IN (
			SELECT post_id
			FROM post_reactions 
			WHERE user_id = $1
		)
		GROUP BY p.id 
	`

	return getAllPostsByQuery(r.DB, query, userID)
}

func (r *postRepository) GetAllCommentedPosts(userID int) (*[]entity.PostEntity, error) {
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
		LEFT JOIN comments cm ON p.id = cm.post_id
		WHERE cm.user_id = $1
		GROUP BY p.id 
	`

	return getAllPostsByQuery(r.DB, query, userID)
}

func (r *postRepository) Exists(postID int) (bool, error) {
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

func (r *postRepository) Delete(postID int, userID int) error {
	query := `
		DELETE FROM posts
		WHERE id = $1 AND user_id = $2
	`

	_, err := r.DB.Exec(query, postID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.ErrForbiddenAccess
		}
		return err
	}

	return nil
}

func (r *postRepository) DeleteByPrivileged(postID int) error {
	query := `
		DELETE FROM posts
		WHERE id = $1
	`

	_, err := r.DB.Exec(query, postID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.ErrPostNotFound
		}
		return err
	}

	return nil
}

func (r *postRepository) GetAuthorID(postID int) (int, error) {
	query := `
		SELECT user_id
		FROM posts
		WHERE id = $1
	`

	var userID int

	err := r.DB.QueryRow(query, postID).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, entity.ErrPostNotFound
		}
		return 0, err
	}

	return userID, nil
}

func (r *postRepository) Update(p entity.PostCreateForm, tagIDs []int, deleteImage bool) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	posts := `
		UPDATE posts
		SET title = $1, content = $2
		WHERE id = $3
	`

	_, err = tx.Exec(posts, p.Title, p.Content, p.ID)
	if err != nil {
		return err
	}

	post_tags_delete := `
		DELETE FROM posts_tags
		WHERE post_id = $1
	`
	_, err = tx.Exec(post_tags_delete, p.ID)
	if err != nil {
		return err
	}

	posts_tags := `
		INSERT INTO posts_tags (post_id, tag_id, created_at)
		VALUES ($1, $2, datetime('now', 'localtime'))
	`
	for _, tagID := range tagIDs {
		_, err := tx.Exec(posts_tags, p.ID, tagID)
		if err != nil {
			return err
		}
	}

	if deleteImage {
		images_delete := `
			DELETE FROM images
			WHERE post_id = $1 
		`
		_, err := tx.Exec(images_delete, p.ID)
		if err != nil {
			return err
		}
	}

	if p.ImageName != "" {
		images := `
			INSERT OR REPLACE INTO images (name, post_id)
			VALUES ($1, $2)
		`
		_, err = tx.Exec(images, p.ImageName, p.ID)
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
