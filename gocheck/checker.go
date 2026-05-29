package gocheck

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"sort"

	"shanhu.io/g/dags"
	"shanhu.io/std/lexing"
)

type checker struct {
	fset  *token.FileSet
	files []*ast.File
	info  *types.Info
	pkg   *types.Package
}

func (c *checker) depGraph() (*dags.Graph, error) {
	depsMap := make(map[token.Pos]map[token.Pos]bool)
	for _, f := range c.files {
		depsMap[filePos(c.fset, f.Pos())] = make(map[token.Pos]bool)
	}

	for use, obj := range c.info.Uses {
		if obj.Pkg() != c.pkg {
			continue // ignore inter-pkg refs
		}

		fused := filePos(c.fset, use.NamePos)
		fdef := filePos(c.fset, obj.Pos())

		if fused == fdef {
			continue
		}

		if _, found := depsMap[fdef]; !found {
			path := c.pkg.Path()
			panic(fmt.Errorf("%s not found in %s", use.Name, path))
		}
		depsMap[fdef][fused] = true
	}

	ret := make(map[string][]string)
	for f, deps := range depsMap {
		var lst []string
		for dep := range deps {
			lst = append(lst, trimBase(filename(c.fset, dep)))
		}
		sort.Strings(lst)
		ret[trimBase(filename(c.fset, f))] = lst
	}
	return &dags.Graph{Nodes: ret}, nil
}

func (c *checker) checkRect(h, w int) []*lexing.Error {
	names := listFileNames(c.fset, c.files)
	return CheckRect(names, h, w)
}

func (c *checker) checkAll(h, w int) []*lexing.Error {
	g, err := c.depGraph()
	if err != nil {
		return lexing.SingleErr(err)
	}

	if err := dags.CheckDAG(g); err != nil {
		return lexing.SingleErr(err)
	}

	names := listFileNames(c.fset, c.files)
	if errs := CheckRect(names, h, w); errs != nil {
		return errs
	}

	return CheckLineComment(c.fset, c.files)
}
