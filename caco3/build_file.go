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
	"path/filepath"

	"shanhu.io/pub/jsonx"
	"shanhu.io/pub/lexing"
	"shanhu.io/pub/osutil"
)

const buildFileName = "BUILD.caco3"

func makeBuildFileNode(t string) interface{} {
	switch t {
	case ruleFileSet:
		return new(FileSet)
	case ruleBundle:
		return new(Bundle)
	case ruleDockerPull:
		return new(DockerPull)
	case ruleDockerBuild:
		return new(DockerBuild)
	case ruleDockerRun:
		return new(DockerRun)
	case ruleDownload:
		return new(Download)
	}
	return nil
}

func readBuildFile(env *env, p string) ([]*buildNode, []*lexing.Error) {
	var fp string
	if p == "" {
		fp = filepath.Join(env.rootDir, buildFileName)
	} else {
		fp = env.src(p, buildFileName)
	}

	if ok, err := osutil.IsRegular(fp); err != nil {
		return nil, lexing.SingleErr(err)
	} else if !ok {
		return nil, nil // No build file present.
	}

	rules, errs := jsonx.ReadSeriesFile(fp, makeBuildFileNode)
	if errs != nil {
		return nil, errs
	}

	var nodes []*buildNode

	errList := lexing.NewErrorList()

	for _, r := range rules {
		node := &buildNode{
			typ:      nodeRule,
			pos:      r.Pos,
			ruleType: r.Type,
		}

		switch v := r.V.(type) {
		case *FileSet:
			fset, err := newFileSet(env, p, v)
			if err != nil {
				errList.Add(&lexing.Error{Pos: r.Pos, Err: err})
				continue
			}
			node.rule = fset
		case *DockerPull:
			dp, err := newDockerPull(env, p, v)
			if err != nil {
				errList.Add(&lexing.Error{Pos: r.Pos, Err: err})
				continue
			}
			node.rule = dp
		case *DockerBuild:
			db, err := newDockerBuild(env, p, v)
			if err != nil {
				errList.Add(&lexing.Error{Pos: r.Pos, Err: err})
				continue
			}
			node.rule = db
		case *DockerRun:
			node.rule = newDockerRun(env, p, v)
		case *Download:
			d, err := newDownload(env, p, v)
			if err != nil {
				errList.Add(&lexing.Error{Pos: r.Pos, Err: err})
				continue
			}
			node.rule = d
		case *Bundle:
			node.rule = newBundle(env, p, v)
		case *SubBuilds:
			node.sub = newSubBuilds(env, p, v)
			node.typ = nodeSub
		default:
			errList.Errorf(r.Pos, "unknown type: %q", r.Type)
			continue
		}

		if node.rule != nil {
			meta, err := node.rule.meta(env)
			if err != nil {
				errList.Errorf(r.Pos, "fail to get rule meta")
			}
			node.ruleMeta = meta
			node.name = meta.name
			node.deps = meta.deps
		}

		if node.name == p || node.name == "" {
			errList.Errorf(r.Pos, "rule has no name")
			continue
		}
		nodes = append(nodes, node)
	}

	if errs := errList.Errs(); errs != nil {
		return nil, errs
	}
	return nodes, nil
}
