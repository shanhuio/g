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

package gocheck

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"
)

func filePos(fset *token.FileSet, p token.Pos) token.Pos {
	return fset.File(p).Pos(0)
}

func filename(fset *token.FileSet, p token.Pos) string {
	return fset.Position(p).Filename
}

func trimBase(name string) string {
	base := filepath.Base(name)
	return strings.TrimSuffix(base, ".go")
}

func listFileNames(fset *token.FileSet, files []*ast.File) []string {
	var names []string
	for _, f := range files {
		tokFile := fset.File(f.Pos())
		name := tokFile.Name()
		names = append(names, name)
	}

	return names
}
