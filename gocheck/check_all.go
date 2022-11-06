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
	"go/build"
	"go/token"

	"golang.org/x/tools/go/packages"
	"shanhu.io/pub/errcode"
	"shanhu.io/pub/lexing"
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
