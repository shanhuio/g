package osutil

import (
	"os"

	"shanhu.io/g/errcode"
)

// CheckPrivateFile checks if the file is of the right permission bits.
func CheckPrivateFile(f string) error {
	info, err := os.Stat(f)
	if err != nil {
		return err
	}
	mod := info.Mode() & 0777
	if mod != 0600 {
		return errcode.InvalidArgf(
			"file %q is not of perm 0600 but %#o", f, mod,
		)
	}
	return err
}

// ReadPrivateFile reads the confent of a file. The file must be mode 0600.
func ReadPrivateFile(f string) ([]byte, error) {
	if err := CheckPrivateFile(f); err != nil {
		return nil, err
	}
	return os.ReadFile(f)
}
