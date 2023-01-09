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

package gocheck

import (
	"go/ast"
	"go/scanner"
	"go/token"
	"io/ioutil"
	"strings"

	"shanhu.io/pub/lexing"
)

func validLineCommentContent(s string) bool {
	if s == "" {
		return true
	}
	if strings.HasPrefix(s, " ") {
		return true
	}
	if strings.HasPrefix(s, "\t") {
		return true
	}
	if strings.HasPrefix(s, "go:build ") {
		return true
	}
	return false
}

func toLexingPos(p token.Position) *lexing.Pos {
	return &lexing.Pos{
		File: p.Filename,
		Line: p.Line,
		Col:  p.Column,
	}
}

func tokenPos(fset *token.FileSet, pos token.Pos) *lexing.Pos {
	return toLexingPos(fset.Position(pos))
}

// CheckLineComment checks the format of line comments.
func CheckLineComment(
	fset *token.FileSet, files []*ast.File,
) []*lexing.Error {
	errs := lexing.NewErrorList()
	errHandler := func(pos token.Position, msg string) {
		errs.Errorf(toLexingPos(pos), "%s", msg)
	}
	for _, f := range files {
		tokFile := fset.File(f.Pos())
		s := new(scanner.Scanner)

		bs, err := ioutil.ReadFile(tokFile.Name())
		if lexing.LogError(errs, err) {
			continue
		}

		s.Init(tokFile, bs, errHandler, scanner.ScanComments)
		for {
			pos, tok, lit := s.Scan()
			if tok == token.EOF {
				break
			}
			if tok == token.COMMENT && strings.HasPrefix(lit, "//") {
				if !validLineCommentContent(lit[2:]) {
					errs.Errorf(
						tokenPos(fset, pos),
						"please add a space to comment %q", lit,
					)
				}
			}
		}
	}
	return errs.Errs()
}
