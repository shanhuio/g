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
