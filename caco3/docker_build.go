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
	"log"
	"os"
	"path"
	"sort"
	"strings"

	"shanhu.io/pub/dock"
	"shanhu.io/pub/errcode"
	"shanhu.io/pub/jsonutil"
	"shanhu.io/pub/strutil"
	"shanhu.io/pub/tarutil"
)

type dockerBuild struct {
	name           string
	rule           *DockerBuild
	fromRuleSums   []string
	dockerfilePath string
	inputs         []string
	archInputs     []string
	prefixDir      string
	repoTag        string
	args           map[string]string
	out            string
	tarOut         string
}

func newDockerBuild(env *env, p string, r *DockerBuild) (
	*dockerBuild, error,
) {
	name := makeRelPath(p, r.Name)

	var f string
	if r.Dockerfile == "" {
		f = path.Join(name, "Dockerfile")
	} else {
		f = makePath(p, r.Dockerfile)
	}

	var fromRuleSums []string
	if len(r.From) > 0 {
		for _, from := range r.From {
			rp := makePath(p, from)
			fromRuleSums = append(fromRuleSums, dockerSumOut(rp))
		}
	}

	repoTag, err := env.nameToRepoTag(name)
	if err != nil {
		return nil, errcode.Annotate(err, "invalid name for docker build")
	}

	args := makeDockerVars(r.Args)

	inputMap := make(map[string]bool)
	for _, input := range r.Input {
		inputMap[makePath(p, input)] = true
	}
	archInputMap := make(map[string]bool)
	for _, input := range r.ArchiveInput {
		archInputMap[makePath(p, input)] = true
	}

	prefixDir := r.PrefixDir
	if prefixDir == "." {
		prefixDir = p
	}

	var tarOut string
	if r.OutputTar {
		tarOut = dockerTarOut(name)
	}

	return &dockerBuild{
		name:           name,
		rule:           r,
		dockerfilePath: f,
		fromRuleSums:   fromRuleSums,
		inputs:         strutil.SortedList(inputMap),
		archInputs:     strutil.SortedList(archInputMap),
		prefixDir:      prefixDir,
		repoTag:        repoTag,
		args:           args,
		out:            dockerSumOut(name),
		tarOut:         tarOut,
	}, nil
}

func (b *dockerBuild) meta(env *env) (*buildRuleMeta, error) {
	dat := struct {
		Dockerfile string            // Know which one is the Dockerfile
		Args       map[string]string `json:",omitempty"`
		PrefixDir  string            `json:",omitempty"`
		OutputTar  bool              `json:",omitempty"`
	}{
		Dockerfile: b.dockerfilePath,
		Args:       b.args,
		PrefixDir:  b.prefixDir,
		OutputTar:  b.rule.OutputTar,
	}

	digest, err := makeDigest(ruleDockerBuild, b.name, &dat)
	if err != nil {
		return nil, errcode.Annotate(err, "digest")
	}

	var deps []string
	deps = append(deps, b.dockerfilePath)
	deps = append(deps, b.fromRuleSums...)
	deps = append(deps, b.inputs...)
	deps = append(deps, b.archInputs...)

	outs := []string{b.out}
	if b.tarOut != "" {
		outs = append(outs, b.tarOut)
	}
	return &buildRuleMeta{
		name:      b.name,
		deps:      strutil.SortedList(strutil.MakeSet(deps)),
		outs:      outs,
		dockerOut: true,
		digest:    digest,
	}, nil
}

func (b *dockerBuild) build(env *env, opts *buildOpts) error {
	dockerfileBytes, err := os.ReadFile(env.src(b.dockerfilePath))
	if err != nil {
		return errcode.Annotate(err, "read Dockerfile")
	}
	df := string(dockerfileBytes)

	ts := dock.NewTarStream(df)
	files := make(map[string]string)

	for _, in := range b.inputs {
		switch typ := env.nodeType(in); typ {
		case "":
			return errcode.Internalf("file %q not found", in)
		case nodeSrc:
			files[in] = env.src(in)
		case nodeOut:
			files[in] = env.out(in)
		case nodeRule:
			fileSet, err := referenceFileSetOut(env, in)
			if err != nil {
				return errcode.Annotatef(err, "input %q", in)
			}
			fileSetFile := env.out(fileSet)
			var list []*fileStat
			if err := jsonutil.ReadFile(fileSetFile, &list); err != nil {
				return errcode.Annotatef(err, "read file set %q", in)
			}
			for _, f := range list {
				var fp string
				switch f.Type {
				case fileTypeSrc:
					fp = env.src(f.Name)
				case fileTypeOut:
					fp = env.out(f.Name)
				default:
					return errcode.Internalf(
						"invalid file type %q of %q ini set %q",
						f.Type, f.Name, in,
					)
				}
				files[f.Name] = fp
			}
		default:
			return errcode.Internalf("unknown type %q", typ)
		}
	}

	var names []string
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)

	prefixDir := b.prefixDir
	if prefixDir != "" && !strings.HasPrefix(prefixDir, "/") {
		prefixDir = prefixDir + "/"
	}

	for _, name := range names {
		tarName := name
		if prefixDir != "" {
			if !strings.HasPrefix(name, prefixDir) {
				continue
			}
			tarName = strings.TrimPrefix(name, prefixDir)
		}

		f := files[name]
		stat, err := os.Stat(f)
		if err != nil {
			return errcode.Annotatef(err, "stat file %q", name)
		}
		mode := stat.Mode()
		if !mode.IsRegular() {
			return errcode.Internalf("%q is not a regular file", name)
		}
		ts.AddFile(tarName, tarutil.ModeMeta(int64(mode)&0777), f)
	}

	for _, ar := range b.archInputs {
		var fp string
		switch typ := env.nodeType(ar); typ {
		case nodeSrc:
			fp = env.src(ar)
		case nodeOut:
			fp = env.out(ar)
		default:
			return errcode.Internalf("unknown type %q", typ)
		}
		base := path.Base(ar)
		dir := path.Dir(ar)
		if dir == "." {
			dir = ""
		}
		if strings.HasSuffix(base, ".zip") {
			ts.AddZipFile(dir, fp)
		} else {
			return errcode.InvalidArgf("unknown archive type %q", base)
		}
	}

	repo, tag := parseRepoTag(b.repoTag)
	rt := repoTag(repo, tag)

	config := &dock.BuildConfig{
		Files:    ts,
		Args:     b.args,
		UseCache: true, // TODO(h8liu): read from option.
	}
	if err := dock.BuildImageConfig(env.dock, rt, config); err != nil {
		return err
	}

	info, err := dock.InspectImage(env.dock, rt)
	if err != nil {
		return errcode.Annotate(err, "inspect built image")
	}

	sum := newDockerSum(repo, tag, info.ID)

	out, err := env.prepareOut(b.out)
	if err != nil {
		return errcode.Annotate(err, "prepare sum output")
	}
	if err := jsonutil.WriteFile(out, sum); err != nil {
		return errcode.Annotate(err, "write image sum")
	}

	if b.tarOut != "" {
		log.Printf("Saving %s", b.tarOut)
		out, err := env.prepareOut(b.tarOut)
		if err != nil {
			return errcode.Annotate(err, "prepare tar output")
		}
		if err := dock.SaveImageGz(env.dock, sum.ID, out); err != nil {
			return errcode.Annotate(err, "save image as tar")
		}
	}

	return nil
}
