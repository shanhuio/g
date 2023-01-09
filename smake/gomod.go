// Copyright (C) 2023  Shanhu Tech Inc.
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

package smake

import (
	"path/filepath"

	"shanhu.io/pub/errcode"
	"shanhu.io/pub/osutil"
)

func findGoModuleRoot(d string) (string, error) {
	for {
		modFile := filepath.Join(d, "go.mod")
		ok, err := osutil.IsRegular(modFile)
		if err != nil {
			return "", errcode.Annotate(err, "check go.mod file")
		}
		if !ok {
			if d == "/" {
				break
			}
			d = filepath.Dir(d)
			if d == "" {
				break
			}
			continue
		}
		return d, nil
	}

	return "", errcode.NotFoundf("go module not found for dir %q", d)
}
