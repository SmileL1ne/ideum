package post

import (
	"fmt"
	"forum/internal/entity"
	"forum/internal/validator"
	"strings"
)

const (
	maxTitleLen   = 100
	maxContentLen = 5000
)

func IsRightPost(p *entity.PostCreateForm) bool {
	p.CheckField(validator.NotBlank(p.Title), "title", "This field cannot be blank")
	p.CheckField(validator.MaxChar(p.Title, maxTitleLen), "title", fmt.Sprintf("Maximum characters length exceeded - %d", maxTitleLen))
	p.CheckField(validator.NotBlank(p.Content), "content", "This field cannot be blank")
	p.CheckField(validator.MaxChar(p.Content, 5000), "content", fmt.Sprintf("Maximum characters length exceeded - %d", maxContentLen))

	return p.Valid()
}

func ConvertEntitiesToViews(posts *[]entity.PostEntity) (*[]entity.PostView, error) {
	// Convert received PostEntity's to PostView's
	var pViews []entity.PostView
	for _, p := range *posts {
		tags := ConvertToStrArr(p.PostTags)
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

func ConvertToStrArr(tagsStr string) []string {
	if tagsStr == "" {
		return []string{}
	}
	return strings.Split(tagsStr, ", ")
}
