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
