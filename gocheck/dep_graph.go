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
	"go/build"

	"shanhu.io/gcimporter"
	"shanhu.io/pub/dags"
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
