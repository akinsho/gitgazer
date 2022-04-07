package main

import (
	"strings"
	"time"
)

func truncateText(s string, max int) string {
	if max > len(s) {
		return s
	}
	return s[:strings.LastIndex(s[:max], " ")] + "..."
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
