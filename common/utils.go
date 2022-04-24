package common

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"unicode/utf8"
)

// OpenURL opens a URL in the default browser of the user.
// This function is intended to be cross-platform.
// @see: https://gist.github.com/nanmu42/4fbaf26c771da58095fa7a9f14f23d27
func OpenURL(url string) (err error) {
	if url == "" {
		return
	}
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return
}

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
