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
	"log"
	"path"
	"sort"
	"strings"

	"shanhu.io/g/dock"
	"shanhu.io/g/errcode"
	"shanhu.io/g/strutil"
	"shanhu.io/g/tarutil"
)

type dockerRun struct {
	name    string
	rule    *DockerRun
	image   string
	ins     map[string]string
	archIns map[string]string
	deps    []string
	outs    []string
	outMap  map[string]string
	envs    map[string]string
}

func newDockerRun(env *env, p string, r *DockerRun) *dockerRun {
	name := makeRelPath(p, r.Name)

	image := makePath(p, r.Image)
	var deps []string
	deps = append(deps, dockerSumOut(image))

	depsMap := make(map[string]bool)
	for _, d := range r.Deps {
		depsMap[makePath(p, d)] = true
	}

	ins := make(map[string]string)
	for f, v := range r.Input {
		inPath := makePath(p, f)
		ins[inPath] = v
		depsMap[inPath] = true
	}
	archIns := make(map[string]string)
	for f, v := range r.ArchiveInput {
		inPath := makePath(p, f)
		archIns[inPath] = v
		depsMap[inPath] = true
	}
	deps = append(deps, strutil.SortedList(depsMap)...)

	var outs []string
	outMap := make(map[string]string)
	for f, v := range r.Output {
		outPath := makeRelPath(p, f)
		outs = append(outs, outPath)
		outMap[outPath] = v
	}
	outs = strutil.SortedList(strutil.MakeSet(outs))

	return &dockerRun{
		name:    name,
		rule:    r,
		image:   image,
		ins:     ins,
		archIns: archIns,
		deps:    deps,
		outs:    outs,
		outMap:  outMap,
		envs:    makeDockerVars(r.Envs),
	}
}

func (r *dockerRun) meta(env *env) (*buildRuleMeta, error) {
	dat := struct {
		Rule *DockerRun
		Envs map[string]string `json:",omitempty"`
	}{
		Rule: r.rule,
		Envs: r.envs,
	}
	digest, err := makeDigest(ruleDockerRun, r.name, &dat)
	if err != nil {
		return nil, errcode.Annotate(err, "digest")
	}

	return &buildRuleMeta{
		name:   r.name,
		outs:   r.outs,
		deps:   r.deps,
		digest: digest,
	}, nil
}

func (r *dockerRun) build(env *env, opts *buildOpts) error {
	contConfig := &dock.ContConfig{
		Cmd:     r.rule.Command,
		WorkDir: r.rule.WorkDir,
		Env:     r.envs,
	}

	if m := r.rule.MountWorkspace; m != "" {
		contConfig.Mounts = append(contConfig.Mounts, &dock.ContMount{
			Host:     env.rootDir,
			Cont:     m,
			ReadOnly: true,
		})
	}

	img, err := env.nameToRepoTag(r.image)
	if err != nil {
		return errcode.Annotate(err, "map image name")
	}

	c := env.dock

	cont, err := dock.CreateCont(c, img, contConfig)
	if err != nil {
		return errcode.Annotate(err, "create container")
	}
	defer cont.Drop()

	if len(r.ins)+len(r.archIns) > 0 {
		ts := tarutil.NewStream()

		var ins []string
		for in := range r.ins {
			ins = append(ins, in)
		}
		sort.Strings(ins)

		for _, in := range ins {
			var f string
			switch typ := env.nodeType(in); typ {
			case "":
				return errcode.Internalf("input %q not found", in)
			case nodeSrc:
				f = env.src(in)
			case nodeOut:
				f = env.out(in)
			default:
				return errcode.Internalf("unknown type %q", typ)
			}

			dest := r.ins[in]
			ts.AddFile(dest, new(tarutil.Meta), f)
		}

		var archIns []string
		for in := range r.archIns {
			archIns = append(archIns, in)
		}
		sort.Strings(archIns)

		for _, in := range archIns {
			var f string
			switch typ := env.nodeType(in); typ {
			case "":
				return errcode.Internalf("archive input %q not found", in)
			case nodeSrc:
				f = env.src(in)
			case nodeOut:
				f = env.out(in)
			default:
				return errcode.Internalf("unknown type %q", typ)
			}
			dest := r.archIns[in]
			base := path.Base(in)
			if strings.HasSuffix(base, ".zip") {
				ts.AddZipFile(dest, f)
			} else {
				return errcode.InvalidArgf("unknown archive type %q", base)
			}
		}

		if err := dock.CopyInTarStream(cont, ts, "/"); err != nil {
			return errcode.Annotate(err, "copy input")
		}
	}

	if err := cont.Start(); err != nil {
		return errcode.Annotate(err, "start container")
	}
	if err := cont.FollowLogs(opts.log); err != nil {
		return errcode.Annotate(err, "stream logs")
	}

	status, err := cont.Wait(dock.NotRunning)
	if err != nil {
		return errcode.Annotate(err, "wait container")
	}
	for _, out := range r.outs {
		from := r.outMap[out]
		to := out

		f, err := env.prepareOut(to)
		if err != nil {
			return errcode.Annotatef(err, "prepare output: %s", to)
		}

		if err := cont.CopyOutFile(from, f); err != nil {
			if status == 0 {
				return errcode.Annotatef(err, "copy %s", to)
			}
			log.Printf("copy %s: %s", to, err)
		}
	}

	if status != 0 {
		return errcode.Internalf("exit with %d", status)
	}

	return nil
}
