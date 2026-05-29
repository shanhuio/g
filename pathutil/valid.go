package pathutil

import (
	"strings"

	"shanhu.io/std/lexing"
)

// ValidPathRune checks if r is a valid rune in the path.
// Valid runes contains a-z, A-Z, 0-9, '_' and '.'
func ValidPathRune(r rune) bool {
	if r == '_' || r == '.' {
		return true
	}
	if r >= 'a' && r <= 'z' {
		return true
	}
	return lexing.IsDigit(r)
}

// ValidPath checks if p is a valid absolute path
func ValidPath(p string) bool {
	if !strings.HasPrefix(p, "/") {
		return false
	}
	if p == "/" {
		return true
	}

	p = strings.TrimPrefix(p, "/")
	subs := strings.Split(p, "/")
	for _, s := range subs {
		if s == "" {
			return false
		}

		for _, r := range s {
			if !ValidPathRune(r) {
				return false
			}
		}
	}
	return true
}
