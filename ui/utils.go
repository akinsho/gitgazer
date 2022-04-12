package ui

import (
	"strings"
	"time"
	"unicode/utf8"
)

func pad(str string, size int) string {
	padding := strings.Repeat(" ", size)
	return padding + str + padding
}

func truncateText(str string, max int, ellipsis bool) string {
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

func throttle(f func(), d time.Duration) func() {
	shouldWait := false
	return func() {
		if !shouldWait {
			f()
			shouldWait = true
			go func() {
				<-time.After(d)
				shouldWait = false
			}()
		}
	}
}