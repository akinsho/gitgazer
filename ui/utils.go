package ui

import (
	"time"
	"unicode/utf8"
)

func truncateText(str string, max int) string {
	if len(str) <= 0 {
		return ""
	}
	if max >= len(str) {
		return str
	}

	if utf8.RuneCountInString(str) < max {
		return str
	}

	return string([]rune(str)[:max]) + "…"
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
