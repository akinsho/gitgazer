package common

import (
	"strings"
	"unicode/utf8"
)

func Pad(str string, size int) string {
	padding := strings.Repeat(" ", size)
	return padding + str + padding
}

func TruncateText(str string, max int, ellipsis bool) string {
	suffix := ""
	if ellipsis {
		suffix = "…"
	}
	if len(str) <= 0 {
		return ""
	}
	if max >= len(str) {
		return str
	}

	if utf8.RuneCountInString(str) < max {
		return str
	}

	return string([]rune(str)[:max]) + suffix
}

func RemoveIndex[T any](s []T, index int) []T {
	if index < 0 || index >= len(s)-1 {
		return s
	}
	ret := make([]T, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}
