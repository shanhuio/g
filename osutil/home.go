package osutil

import (
	"os"
	"path"
	"path/filepath"

	"shanhu.io/std/errcode"
)

// Home is a directory for referecing files under a directory.
type Home struct {
	dir string
}

// NewHome creates a new home directory.
func NewHome(dir string) (*Home, error) {
	if dir == "" {
		wd, err := os.Getwd()
		if err != nil {
			return nil, errcode.Annotate(err, "get working dir")
		}
		dir = wd
	} else {
		abs, err := filepath.Abs(dir)
		if err != nil {
			return nil, errcode.Annotate(err, "get absolute dir")
		}
		dir = abs
	}

	return &Home{dir: dir}, nil
}

// FilePath returns a sub path under the home directory. p is in URL path, but
// the returned value is in filepath format, in OS's filepath separators.
func (h *Home) FilePath(p string) string {
	if p == "" {
		return h.dir
	}
	return filepath.Join(h.dir, filepath.FromSlash(p))
}

// Dir returns the base directory, it is always in its absolute form.
func (h *Home) Dir() string { return h.dir }

// Var returns a file path under dir "var/"
func (h *Home) Var(p string) string {
	return h.FilePath(path.Join("var", p))
}

// Etc returns a file path under dir "etc/"
func (h *Home) Etc(p string) string {
	return h.FilePath(path.Join("etc", p))
}

// Lib returns a file path under dir "lib/"
func (h *Home) Lib(p string) string {
	return h.FilePath(path.Join("lib", p))
}

// Tmp returns a file path under dir "tmp/"
func (h *Home) Tmp(p string) string {
	return h.FilePath(path.Join("tmp", p))
}
