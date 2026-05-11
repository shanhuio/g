package pathutil

import (
	"strings"
)

// IsParent checks if path short is a parent of path long.
func IsParent(short, long string) bool {
	if short == long {
		return true
	}
	return strings.HasPrefix(long, short+"/")
}
