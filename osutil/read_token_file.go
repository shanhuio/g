package osutil

import (
	"os"
	"strings"
)

// ReadTokenFile reads a token string from a file.
func ReadTokenFile(f string) (string, error) {
	bs, err := os.ReadFile(f)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bs)), nil
}

// ReadOptionalTokenFile reads an optional token file.
func ReadOptionalTokenFile(f string) (string, bool, error) {
	ret, err := ReadTokenFile(f)
	if err != nil {
		if os.IsNotExist(err) {
			return "", false, nil
		}
		return "", false, err
	}
	return ret, true, nil
}
