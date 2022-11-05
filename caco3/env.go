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

package caco3

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"shanhu.io/pub/dock"
	"shanhu.io/pub/errcode"
)

type env struct {
	dock *dock.Client

	rootDir     string
	workDir     string
	workSrcPath string

	srcDir string
	outDir string

	workspace *Workspace // Lazily loaded.

	nodeType func(name string) string
	ruleType func(name string) string
}

func (e *env) prepareOut(ps ...string) (string, error) {
	p := e.out(ps...)
	dir := filepath.Dir(p)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}
	return p, nil
}

func dirFilePath(dir string, ps ...string) string {
	if len(ps) == 0 {
		return dir
	}
	p := path.Join(ps...)
	return filepath.Join(dir, filepath.FromSlash(p))
}

func (e *env) root(ps ...string) string {
	return dirFilePath(e.rootDir, ps...)
}

func (e *env) out(ps ...string) string {
	return dirFilePath(e.outDir, ps...)
}

func (e *env) src(ps ...string) string {
	return dirFilePath(e.srcDir, ps...)
}

func (e *env) nameToRepoTag(name string) (string, error) {
	parts := strings.Split(name, "/")
	if len(parts) == 0 {
		return "", errcode.InvalidArgf("empty name")
	}
	if len(parts) != 4 {
		return "", errcode.InvalidArgf("invalid name %q", name)
	}

	domain := parts[0]
	project := parts[1]
	dockers := parts[2]
	base := parts[3]

	if dockers != "dockers" && !strings.HasSuffix(dockers, "-dockers") {
		return "", errcode.InvalidArgf("not a docker image name: %q", name)
	}

	if domain == "shanhu.io" {
		domain = "cr.shanhu.io"
	}

	return path.Join(domain, project, base), nil
}
