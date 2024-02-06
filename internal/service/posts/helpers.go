package posts

import (
	"fmt"
	"forum/internal/entity"
	"forum/internal/validator"
)

const (
	maxTitleLen   = 100
	maxContentLen = 5000
)

func isRightPost(p *entity.PostCreateForm) bool {
	p.CheckField(validator.NotBlank(p.Title), "title", "This field cannot be blank")
	p.CheckField(validator.MaxChar(p.Title, maxTitleLen), "title", fmt.Sprintf("Maximum characters length exceeded - %d", maxTitleLen))
	p.CheckField(validator.NotBlank(p.Content), "content", "This field cannot be blank")
	p.CheckField(validator.MaxChar(p.Content, 5000), "content", fmt.Sprintf("Maximum characters length exceeded - %d", maxContentLen))

	return p.Valid()
}
