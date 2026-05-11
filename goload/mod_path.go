package goload

import (
	"strconv"
	"strings"
)

func isValidModPath(p, modPath string) bool {
	if modPath == p {
		return false
	}

	prefix := p + "/v"
	if !strings.HasPrefix(modPath, prefix) {
		return false
	}

	ver := strings.TrimPrefix(modPath, prefix)
	if _, err := strconv.Atoi(ver); err != nil {
		return false
	}

	return true
}
