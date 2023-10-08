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

package caco3

import (
	"io/fs"
	"os"
	"sort"

	"shanhu.io/g/errcode"
	"shanhu.io/g/lexing"
)

type loader struct {
	env *env

	// All parsed and registered build nodes.
	nodes map[string]*buildNode

	// All loaded build nodes. A loaded node always has its dependencies
	// loaded.
	loaded map[string]*buildNode

	tracer *loadTracer

	errList *lexing.ErrorList
}

func newLoader(env *env) *loader {
	return &loader{
		env:     env,
		loaded:  make(map[string]*buildNode),
		nodes:   make(map[string]*buildNode),
		tracer:  newLoadTracer(),
		errList: lexing.NewErrorList(),
	}
}

func (l *loader) register(n *buildNode) {
	if n.name == "" {
		l.errList.Errorf(n.pos, "node name is empty")
		return
	}
	if p, ok := l.nodes[n.name]; ok {
		l.errList.Errorf(n.pos, "node with name %q redeclared", n.name)
		if p.pos != nil {
			l.errList.Errorf(p.pos, "  previously defined here")
		}
		return
	}
	l.nodes[n.name] = n
}

// load all names that is referenced at pos.
func (l *loader) load(names []string, pos *lexing.Pos) []*buildNode {
	var nodes []*buildNode
	for _, name := range names {
		n := l.load1(name, pos)
		nodes = append(nodes, n)
	}
	return nodes
}

func (l *loader) load1(name string, pos *lexing.Pos) *buildNode {
	if !l.tracer.push(name) {
		l.errList.Errorf(
			pos, "has circular dependency: %q", l.tracer.stack(),
		)
		return nil
	}
	defer l.tracer.pop()

	if n, ok := l.loaded[name]; ok {
		return n // already loaded
	}

	n, ok := l.nodes[name]
	if ok { // Registered but not loaded yet
		l.load(n.deps, pos) // Load its dependencies.
		l.loaded[name] = n  // Add into loaded map.
		return n
	}

	// Auto register and load source files.
	f := l.env.src(name)
	stat, err := os.Lstat(f)
	if err != nil {
		l.errList.Errorf(pos, "stat %q: %s", f, err)
		return nil
	}

	mode := stat.Mode()
	if mode.IsRegular() || mode.Type() == fs.ModeSymlink {
		n := &buildNode{
			name: name,
			typ:  nodeSrc,
		}
		l.register(n)
		l.loaded[name] = n
		return n
	}

	l.errList.Errorf(pos, "cannot resolve %q", name)
	return nil
}

func (l *loader) registerOuts(
	rule string, names []string, pos *lexing.Pos,
) {
	if len(names) == 0 {
		return
	}

	deps := []string{rule}
	for _, name := range names {
		n := &buildNode{
			name: name,
			typ:  nodeOut,
			deps: deps,
			pos:  pos,
		}
		l.register(n)
	}
}

func (l *loader) readBuildFile(p string) {
	subDirMap := make(map[string]bool)

	nodes, errs := readBuildFile(l.env, p)
	l.errList.AddAll(errs)
	for _, n := range nodes {
		if n.typ == nodeSub {
			for _, d := range n.sub.Dirs() {
				subDirMap[d] = true
			}
			continue
		}

		l.register(n)

		if n.typ == nodeRule {
			l.registerOuts(n.name, n.ruleMeta.outs, n.pos)
		}
	}

	var subDirs []string
	for subDir := range subDirMap {
		subDirs = append(subDirs, subDir)
	}
	sort.Strings(subDirs)

	for _, d := range subDirs {
		l.readBuildFile(d)
	}
}

func (l *loader) Errs() []*lexing.Error {
	return l.errList.Errs()
}

func loadNodes(env *env, names []string) (
	[]*buildNode, map[string]*buildNode, []*lexing.Error,
) {
	l := newLoader(env)

	repoMap := env.workspace.RepoMap
	if repoMap == nil || len(repoMap.Src) == 0 {
		err := errcode.InvalidArgf("repo map missing")
		return nil, nil, lexing.SingleErr(err)
	}
	var dirs []string
	for dir := range repoMap.Src {
		dirs = append(dirs, dir)
	}
	sort.Strings(dirs)

	for _, dir := range dirs {
		l.readBuildFile(dir)
	}

	if errs := l.Errs(); errs != nil {
		return nil, nil, errs
	}

	nodes := l.load(names, nil)
	if errs := l.Errs(); errs != nil {
		return nil, nil, errs
	}

	return nodes, l.loaded, nil
}
