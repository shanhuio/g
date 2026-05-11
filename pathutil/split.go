package pathutil

import (
	"fmt"
	"strings"
)

// Split splits the package name into parts.
func Split(path string) ([]string, error) {
	if path == "" {
		return nil, fmt.Errorf("empty path")
	}

	parts := strings.Split(path, "/")
	for _, part := range parts {
		if part == "" {
			return nil, fmt.Errorf("invalid path: %q", path)
		}
	}
	return parts, nil
}
