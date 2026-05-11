// Package gomod provides simple go.mod file parsing.
//
// TODO(h8liu): simply call golang.org/x/mod/modfile
package gomod

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"strconv"
	"strings"
)

// File is a parsed go.mod file.
type File struct {
	Name string
}

// Parse parses a go.mod file.
func Parse(f string) (*File, error) {
	bs, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}
	name, err := modulePath(bs)
	if err != nil {
		return nil, err
	}
	return &File{Name: name}, nil
}

var errInvalidModFile = errors.New("invalid go.mod file")

// modulePath returns the module path from the gomod file text.
// If it cannot find a module path, it returns an empty string.
// It is tolerant of unrelated problems in the go.mod file.
func modulePath(bs []byte) (string, error) {
	s := bufio.NewScanner(bytes.NewReader(bs))

	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		if !strings.HasPrefix(line, "module") {
			continue
		}

		line = strings.TrimSpace(strings.TrimPrefix(line, "module"))

		// TODO: this is incorrect for quoted module path
		if before, _, found := strings.Cut(line, "//"); found {
			line = strings.TrimSpace(before)
			if line == "" {
				return "", errInvalidModFile
			}
		}

		if line == "" {
			return "", errInvalidModFile
		}
		if line[0] == '"' || line[0] == '`' {
			p, err := strconv.Unquote(line)
			if err != nil || p == "" {
				return "", errInvalidModFile
			}
			return p, nil
		}

		return line, nil
	}

	return "", errInvalidModFile
}
