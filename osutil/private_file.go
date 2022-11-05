// Copyright (C) 2022  Shanhu Tech Inc.
//
// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU Affero General Public License as published by the
// Free Software Foundation, either version 3 of the License, or (at your
// option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.  See the GNU Affero General Public License
// for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package osutil

import (
	"os"

	"shanhu.io/pub/errcode"
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
