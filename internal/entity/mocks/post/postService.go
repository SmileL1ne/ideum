package post

import (
	"errors"
	"forum/internal/entity"
	repo "forum/internal/repository/post"
	service "forum/internal/service/post"
	"strings"
)

type PostServiceMock struct {
	pr repo.IPostRepository
}

func NewPostServiceMock(r repo.IPostRepository) *PostServiceMock {
	return &PostServiceMock{
		pr: r,
	}
}

var _ service.IPostService = (*PostServiceMock)(nil)

func (ps *PostServiceMock) SavePost(*entity.PostCreateForm, int, []string) (int, error) {
	return ps.pr.SavePost(entity.PostCreateForm{}, 0, nil)
}

func (ps *PostServiceMock) GetPost(postID int) (entity.PostView, error) {
	postEntity, err := ps.pr.GetPost(postID)
	if errors.Is(err, entity.ErrNoRecord) {
		return entity.PostView{}, entity.ErrInvalidPostID
	}

	return entity.PostView{
		ID:          postEntity.ID,
		Title:       postEntity.Title,
		Content:     postEntity.Content,
		Username:    postEntity.Username,
		Likes:       postEntity.Likes,
		Dislikes:    postEntity.Dislikes,
		PostTags:    convertToStrArr(postEntity.PostTags),
		CommentsLen: postEntity.CommentsLen,
	}, nil
}

func (ps *PostServiceMock) GetAllPosts() (*[]entity.PostView, error) {
	posts, _ := ps.pr.GetAllPosts()
	return convertEntitiesToViews(posts)
}

func (ps *PostServiceMock) GetAllPostsByTagId(tagID int) (*[]entity.PostView, error) {
	posts, _ := ps.pr.GetAllPostsByTagId(tagID)
	return convertEntitiesToViews(posts)
}

func (ps *PostServiceMock) GetAllPostsByUserId(userID int) (*[]entity.PostView, error) {
	posts, _ := ps.pr.GetAllPostsByUserID(userID)
	return convertEntitiesToViews(posts)
}

func (ps *PostServiceMock) GetAllPostsByUserReaction(userID int) (*[]entity.PostView, error) {
	posts, _ := ps.pr.GetAllPostsByUserReaction(userID)
	return convertEntitiesToViews(posts)
}

func (ps *PostServiceMock) ExistsPost(postID int) (bool, error) {
	return ps.pr.ExistsPost(postID)
}

func convertEntitiesToViews(posts *[]entity.PostEntity) (*[]entity.PostView, error) {
	var pViews []entity.PostView
	for _, p := range *posts {
		tags := convertToStrArr(p.PostTags)
		post := entity.PostView{
			ID:          p.ID,
			Title:       p.Title,
			Content:     p.Content,
			CreatedAt:   p.CreatedAt,
			Username:    p.Username,
			Likes:       p.Likes,
			Dislikes:    p.Dislikes,
			PostTags:    tags,
			CommentsLen: p.CommentsLen,
		}
		pViews = append(pViews, post)
	}

	return &pViews, nil
}

func convertToStrArr(tagsStr string) []string {
	if tagsStr == "" {
		return []string{}
	}
	return strings.Split(tagsStr, ", ")
}
