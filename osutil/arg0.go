package osutil

import (
	"os"
	"path/filepath"
)

// Arg0 returns the first arg, often represents the path of the binary.
func Arg0() string {
	if len(os.Args) == 0 {
		return ""
	}
	return os.Args[0]
}

// Arg0Base returns the base name of the first arg, which often represents the
// name of the binary.
func Arg0Base() string {
	return filepath.Base(Arg0())
}
