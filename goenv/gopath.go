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

package goenv

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

// GOPATH returns GOPATH reading from environment variables.
// If GOPATH is missing it returns $HOME/go.
func GOPATH() (string, error) {
	p := os.Getenv("GOPATH")
	if p != "" {
		lst := filepath.SplitList(p)
		if len(lst) > 1 {
			return "", fmt.Errorf("GOPATH contains multiple folders")
		}
		return p, nil
	}

	u, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(u.HomeDir, "go"), nil
}
