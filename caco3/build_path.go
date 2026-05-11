package caco3

import (
	"path"
	"strings"
)

// makeRelPath makes a path that is under p. It cannot escape p.
func makeRelPath(p, f string) string {
	f = path.Clean(path.Join("/", f))
	return strings.TrimPrefix(path.Join("/", p, f), "/")
}

func makePath(p, f string) string {
	if path.IsAbs(f) {
		return strings.TrimPrefix(path.Clean(f), "/")
	}
	return makeRelPath(p, f)
}
