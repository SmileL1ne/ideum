package post

import (
	"database/sql"
	"errors"
	"forum/internal/entity"
	"sync"
)

type IPostRepository interface {
	SavePost(entity.PostCreateForm, int, []int) (int, error)
	GetPost(int) (entity.PostEntity, error)
	GetAllPosts() (*[]entity.PostEntity, error)
	GetAllPostsByTagId(int) (*[]entity.PostEntity, error)
	GetTagsForEachPost(*[]entity.PostEntity) (*[]entity.PostEntity, error)
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

func (r *postRepository) SavePost(p entity.PostCreateForm, userID int, tagIDs []int) (int, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return 0, err
	}

	query1 := `
		INSERT INTO posts (title, content, user_id, created_at) 
		VALUES ($1, $2, $3, datetime('now', 'utc', '+12 hours'))
	`

	result, err := tx.Exec(query1, p.Title, p.Content, userID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	query2 := `
		INSERT INTO posts_tags (post_id, tag_id, created_at)
		VALUES ($1, $2, datetime('now', 'utc', '+12 hours'))
	`

	for _, tagID := range tagIDs {
		_, err := tx.Exec(query2, postID, tagID)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
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
			SUM(CASE WHEN pr.is_like = false THEN 1 ELSE 0 END) as dislikes_count  
		FROM posts p
		INNER JOIN users u ON p.user_id = u.id
		LEFT JOIN post_reactions pr ON p.id = pr.post_id
		WHERE p.id=$1
		GROUP BY p.id
		`

	var post entity.PostEntity
	if err := r.DB.QueryRow(query, postID).Scan(&post.ID, &post.Title, &post.Content,
		&post.CreatedAt, &post.Username, &post.Likes, &post.Dislikes); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.PostEntity{}, entity.ErrNoRecord
		}
		return entity.PostEntity{}, err
	}

	return post, nil
}

func (r *postRepository) GetAllPosts() (*[]entity.PostEntity, error) {
	query := `
		SELECT p.id, p.title, p.content, p.created_at, u.username, 
			SUM(CASE WHEN pr.is_like = true THEN 1 ELSE 0 END) as likes_count,
			SUM(CASE WHEN pr.is_like = false THEN 1 ELSE 0 END) as dislikes_count
		FROM posts p
		INNER JOIN users u ON p.user_id = u.id
		LEFT JOIN post_reactions pr ON p.id = pr.post_id
		GROUP BY p.id
	`

	return r.getAllPostsByQuery(query)
}

func (r *postRepository) GetAllPostsByTagId(tagID int) (*[]entity.PostEntity, error) {
	query := `
		SELECT p.id, p.title, p.content, p.created_at, u.username, 
			SUM(CASE WHEN pr.is_like = true THEN 1 ELSE 0 END) as likes_count,
			SUM(CASE WHEN pr.is_like = false THEN 1 ELSE 0 END) as dislikes_count
		FROM posts p
		INNER JOIN users u ON p.user_id = u.id
		LEFT JOIN post_reactions pr ON p.id = pr.post_id
		WHERE p.id IN (
			SELECT pt.post_id
			FROM posts_tags pt
			WHERE pt.tag_id = $1
			)
		GROUP BY p.id
	`

	return r.getAllPostsByQuery(query, tagID)
}

func (r *postRepository) getAllPostsByQuery(query string, args ...interface{}) (*[]entity.PostEntity, error) {
	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}

	var posts []entity.PostEntity

	for rows.Next() {
		var post entity.PostEntity
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt,
			&post.Username, &post.Likes, &post.Dislikes); err != nil {

			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &posts, nil
}

func (r *postRepository) GetTagsForEachPost(posts *[]entity.PostEntity) (*[]entity.PostEntity, error) {
	query := `
		SELECT t.id, t.name, t.created_at
		FROM tags t
		LEFT JOIN posts_tags pt ON pt.tag_id = t.id
		WHERE pt.post_id = $1
	`

	var wg sync.WaitGroup
	wg.Add(len(*posts))
	errCh := make(chan error, len(*posts))

	for i := 0; i < len(*posts); i++ {
		go func(p *entity.PostEntity) {
			defer wg.Done()

			var tags []entity.TagEntity
			rows, err := r.DB.Query(query, p.ID)
			if err != nil {
				errCh <- err
				return
			}
			for rows.Next() {
				var tag entity.TagEntity
				if err := rows.Scan(&tag.ID, &tag.Name, &tag.CreatedAt); err != nil {
					errCh <- err
					return
				}
				tags = append(tags, tag)
			}
			p.PostTags = tags
		}(&(*posts)[i])
	}

	wg.Wait()

	select {
	case err := <-errCh:
		return nil, err
	default:
		return posts, nil
	}
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
