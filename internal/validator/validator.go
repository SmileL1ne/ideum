package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

/*
	Validator is package that provides tools for validating
	data.

	Validator struct holds 2 structures that holds field and
	non-field error messages if there any
*/

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

func MinChar(str string, n int) bool {
	return utf8.RuneCountInString(str) >= n
}

func Matches(str string, rx *regexp.Regexp) bool {
	return rx.MatchString(str)
}
