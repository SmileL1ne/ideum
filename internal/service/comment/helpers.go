package comment

import (
	"fmt"
	"forum/internal/entity"
	"forum/internal/validator"
)

const (
	commentMaxLen = 500
)

func IsRightComment(c *entity.CommentCreateForm) bool {
	c.CheckField(validator.NotBlank(c.Content), "commentContent", "This field cannot be blank")
	c.CheckField(validator.MaxChar(c.Content, commentMaxLen), "commentContent", fmt.Sprintf("This cannot be longer than %d characters", commentMaxLen))

	return c.Valid()
}
