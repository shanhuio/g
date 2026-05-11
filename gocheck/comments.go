package gocheck

import (
	"go/ast"
	"go/scanner"
	"go/token"
	"io/ioutil"
	"strings"

	"shanhu.io/g/lexing"
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
