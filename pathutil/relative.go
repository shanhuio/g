package pathutil

import (
	"strings"
)

// DotRelative returns the relative path of full to base.
// The return value starts with a dot if base is a parent of full.
// It returns full unchanged if base is not a parent of full.
func DotRelative(base, full string) string {
	if base == full {
		return "."
	}
	if strings.HasPrefix(full, base+"/") {
		return "./" + strings.TrimPrefix(full, base+"/")
	}
	return full
}

// Relative returns the relative path of full to base.
// It returns the relative path if base is a parent of full.
// It returns a single dot if base and full are the same.
// It retunrs an empty string if base is not a parent of full.
func Relative(base, full string) string {
	if base == full {
		return "."
	}
	if strings.HasPrefix(full, base+"/") {
		return strings.TrimPrefix(full, base+"/")
	}
	return ""
}
