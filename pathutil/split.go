package pathutil

import (
	"fmt"
	"slices"
	"strings"
)

// Split splits the package name into parts.
func Split(path string) ([]string, error) {
	if path == "" {
		return nil, fmt.Errorf("empty path")
	}

	parts := strings.Split(path, "/")
	if slices.Contains(parts, "") {
		return nil, fmt.Errorf("invalid path: %q", path)
	}
	return parts, nil
}
