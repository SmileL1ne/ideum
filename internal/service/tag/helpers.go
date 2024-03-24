package tag

import (
	"strings"
	"unicode/utf8"
)

const (
	tagMinLen = 3
	tagMaxLen = 30
)

func IsValidTag(tag string) bool {
	trimmed := strings.TrimSpace(tag)

	if len(trimmed) == 0 {
		return false
	}

	if utf8.RuneCountInString(trimmed) < tagMinLen || utf8.RuneCountInString(trimmed) > tagMaxLen {
		return false
	}

	return true
}
