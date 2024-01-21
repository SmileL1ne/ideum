package users

import (
	"fmt"
	"forum/internal/entity"
	"forum/internal/service/validator"
	"regexp"
)

const (
	maxUsernameLen = 255
	maxEmailLen    = 255
	minPasswordLen = 8
	maxPasswordLen = 500
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:.[a-zA-Z0-9](?:[a-zA-Z0-9]{0, 61}[a-zA-Z0-9])?)*$")

func isRightUser(u *entity.UserSignupForm) bool {
	u.CheckField(validator.NotBlank(u.Username), "username", "This field cannot be blank")
	u.CheckField(validator.MaxChar(u.Username, maxUsernameLen), "username", fmt.Sprintf("Maximum characters length exceeded - %d", maxUsernameLen))
	u.CheckField(validator.NotBlank(u.Email), "email", "This field cannot be blank")
	u.CheckField(validator.MaxChar(u.Email, maxEmailLen), "email", fmt.Sprintf("Maximum characters length exceeded - %d", maxEmailLen))
	u.CheckField(validator.Matches(u.Email, EmailRX), "email", "Invalid email address")
	u.CheckField(validator.NotBlank(u.Password), "password", "This field cannot be blank")
	u.CheckField(validator.MinChar(u.Password, minPasswordLen), "password", fmt.Sprintf("This field should be %d characters length minimum", minPasswordLen))
	u.CheckField(validator.MaxChar(u.Password, maxPasswordLen), "password", fmt.Sprintf("Maximum characters length exceeded - %d", maxPasswordLen))

	return u.Valid()
}
