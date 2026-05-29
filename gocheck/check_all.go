package gocheck

import (
	"go/build"
	"go/token"

	"golang.org/x/tools/go/packages"
	"shanhu.io/std/errcode"
	"shanhu.io/std/lexing"
)

// ModCheckAll performs all checks on the package.
func ModCheckAll(dir, pkg string, h, w int) []*lexing.Error {
	var loadMode packages.LoadMode
	for _, m := range []packages.LoadMode{
		packages.NeedTypes,
		packages.NeedFiles,
		packages.NeedTypesInfo,
		packages.NeedModule,
		packages.NeedSyntax,
	} {
		loadMode |= m
	}

	fset := token.NewFileSet()
	config := &packages.Config{
		Mode: loadMode,
		Dir:  dir,
		Fset: fset,
	}
	pkgs, err := packages.Load(config, pkg)
	if err != nil {
		return lexing.SingleErr(err)
	}

	if len(pkgs) != 1 {
		err := errcode.Internalf("got %d packages", len(pkgs))
		return lexing.SingleErr(err)
	}

	p := pkgs[0]

	c := &checker{
		fset:  fset,
		files: p.Syntax,
		info:  p.TypesInfo,
		pkg:   p.Types,
	}
	return c.checkAll(h, w)
}

// CheckAll checks everything for a package.
func CheckAll(path string, h, w int) []*lexing.Error {
	l, err := newLoaderPath(&build.Default, path, nil)
	if err != nil {
		return lexing.SingleErr(err)
	}
	c, err := l.checker()
	if err != nil {
		return lexing.SingleErr(err)
	}
	return c.checkAll(h, w)
}
