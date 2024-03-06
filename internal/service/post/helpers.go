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

var types = map[interface{}]struct{}{
	"image/jpeg": {},
	"image/png":  {},
	"image/gif":  {},
	"image/jpg":  {},
}

func IsRightPost(p *entity.PostCreateForm, withImage bool) bool {
	p.CheckField(validator.NotBlank(p.Title), "title", "This field cannot be blank")
	p.CheckField(validator.MaxChar(p.Title, maxTitleLen), "title", fmt.Sprintf("Maximum characters length exceeded - %d", maxTitleLen))
	p.CheckField(validator.NotBlank(p.Content), "content", "This field cannot be blank")
	p.CheckField(validator.MaxChar(p.Content, 5000), "content", fmt.Sprintf("Maximum characters length exceeded - %d", maxContentLen))
	p.CheckField(validator.NotZero(len(p.Tags)), "tags", "At least one tag should be selected")

	if withImage {
		contentType := p.FileHeader.Header.Get("Content-Type")
		p.CheckField(validator.ExistsInSet(contentType, types), "image", "Only these types are allowed - '.jpeg', '.png', '.gif', '.jpg'")
		p.CheckField(validator.LessThan(p.FileHeader.Size, (20<<20)), "image", "Image size is too large")
	}

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
			ImageName:   p.ImageName,
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
