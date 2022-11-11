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

package smake

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"path/filepath"
)

func listFiles(pkg *build.Package) []string {
	var files []string
	files = append(files, pkg.GoFiles...)
	files = append(files, pkg.CgoFiles...)
	files = append(files, pkg.TestGoFiles...)
	return files
}

func listAbsFiles(pkg *build.Package) []string {
	files := listFiles(pkg)
	for i, f := range files {
		files[i] = filepath.Join(pkg.Dir, f)
	}
	return files
}

func fileSourceMap(pkg *relPkg) (map[string][]byte, error) {
	files := listFiles(pkg.pkg)
	fileMap := make(map[string][]byte)

	for _, f := range files {
		path := filepath.Join(pkg.pkg.Dir, f)
		src, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read %q: %s", path, err)
		}
		fileMap[path] = src
	}

	return fileMap, nil
}
