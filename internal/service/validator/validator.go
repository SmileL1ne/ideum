package validator

import (
	"strings"
	"unicode/utf8"
)

type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

func (v *Validator) CheckField(isRight bool, key, message string) {
	if !isRight {
		v.AddFieldError(key, message)
	}
}

func NotBlank(str string) bool {
	return strings.TrimSpace(str) != ""
}

func MaxChar(str string, n int) bool {
	return utf8.RuneCountInString(str) <= n
}
