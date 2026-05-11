package dock

import (
	"path"
)

func contPath(id, method string) string {
	p := path.Join("containers", id)
	if method == "" {
		return p
	}
	return path.Join(p, method)
}

func execPath(id, method string) string {
	return path.Join("exec", id, method)
}
