package gocheck

import (
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"go/types"
	"path/filepath"

	"shanhu.io/gcimporter"
)

type loader struct {
	path string

	ctx      *build.Context
	buildPkg *build.Package

	fset  *token.FileSet
	alias *gcimporter.AliasMap
}

func newLoaderPath(
	ctx *build.Context, path string, alias *gcimporter.AliasMap,
) (*loader, error) {
	if alias != nil {
		path = alias.Map(path)
	}
	pkg, err := ctx.Import(path, "", 0)
	if err != nil {
		return nil, err
	}

	return newLoader(ctx, pkg, alias), nil
}

func newLoader(
	ctx *build.Context, pkg *build.Package, alias *gcimporter.AliasMap,
) *loader {
	fset := token.NewFileSet()
	return &loader{
		ctx:      ctx,
		path:     pkg.ImportPath,
		buildPkg: pkg,
		fset:     fset,
		alias:    alias,
	}
}

func (l *loader) listFiles() ([]*ast.File, error) {
	var srcFiles []string
	srcFiles = append(srcFiles, l.buildPkg.GoFiles...)
	srcFiles = append(srcFiles, l.buildPkg.CgoFiles...)

	var files []*ast.File
	for _, baseName := range srcFiles {
		filename := filepath.Join(l.buildPkg.Dir, baseName)
		f, err := parser.ParseFile(l.fset, filename, nil, 0)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}

	return files, nil
}

func (l *loader) typesCheck(files []*ast.File) (
	*types.Info, *types.Package, error,
) {
	config := &types.Config{
		Importer:    gcimporter.New(l.ctx, l.alias),
		FakeImportC: true,
	}
	info := &types.Info{
		Uses: make(map[*ast.Ident]types.Object),
	}

	typesPkg, err := config.Check(l.path, l.fset, files, info)
	if err != nil {
		return nil, nil, err
	}
	return info, typesPkg, nil
}

func (l *loader) checker() (*checker, error) {
	files, err := l.listFiles()
	if err != nil {
		return nil, err
	}
	info, typesPkg, err := l.typesCheck(files)
	if err != nil {
		return nil, err
	}

	return &checker{
		fset:  l.fset,
		files: files,
		info:  info,
		pkg:   typesPkg,
	}, nil
}
