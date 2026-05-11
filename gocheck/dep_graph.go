package gocheck

import (
	"go/build"

	"shanhu.io/g/dags"
	"shanhu.io/gcimporter"
)

// DepGraph returns the dependency graph for files in a package.
func DepGraph(
	ctx *build.Context, path string, alias *gcimporter.AliasMap,
) (*dags.Graph, error) {
	l, err := newLoaderPath(ctx, path, alias)
	if err != nil {
		return nil, err
	}
	c, err := l.checker()
	if err != nil {
		return nil, err
	}
	return c.depGraph()
}

// DepGraphPkg returns the dependency graph for files in a loaded package.
func DepGraphPkg(
	ctx *build.Context, pkg *build.Package, alias *gcimporter.AliasMap,
) (*dags.Graph, error) {
	l := newLoader(ctx, pkg, alias)
	c, err := l.checker()
	if err != nil {
		return nil, err
	}
	return c.depGraph()
}
