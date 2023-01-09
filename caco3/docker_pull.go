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
	"fmt"
	"log"
	"strings"

	"shanhu.io/pub/dock"
	"shanhu.io/pub/errcode"
	"shanhu.io/pub/jsonutil"
)

type dockerPull struct {
	name    string
	rule    *DockerPull
	repoTag string
	out     string
	tarOut  string
}

func newDockerPull(env *env, p string, r *DockerPull) (*dockerPull, error) {
	name := makeRelPath(p, r.Name)
	repoTag, err := env.nameToRepoTag(name)
	if err != nil {
		return nil, errcode.Annotate(err, "invalid docker pull name")
	}
	pull := &dockerPull{
		name:    name,
		rule:    r,
		repoTag: repoTag,
		out:     dockerSumOut(name),
	}
	if r.OutputTar {
		pull.tarOut = dockerTarOut(name)
	}
	return pull, nil
}

func (p *dockerPull) pull(env *env) (*dockerSum, error) {
	r := p.rule

	repo, tag := parseRepoTag(p.repoTag)
	srcRepo, srcTag := repo, tag

	if r.Pull != "" {
		srcRepo, srcTag = parseRepoTag(r.Pull)
	}

	digest := r.Digest

	from := repoTag(srcRepo, srcTag)
	pullTag := srcTag

	if digest != "" {
		from = fmt.Sprintf("%s@%s", srcRepo, digest)
		pullTag = digest
	}

	if err := dock.PullImage(env.dock, srcRepo, pullTag); err != nil {
		return nil, errcode.Annotate(err, "pull image")
	}
	if err := dock.TagImage(env.dock, from, srcRepo, srcTag); err != nil {
		return nil, errcode.Annotate(err, "tag image as source")
	}
	if !(repo == srcRepo && tag == srcTag) {
		if err := dock.TagImage(env.dock, from, repo, tag); err != nil {
			return nil, errcode.Annotate(err, "re-tag output image")
		}
	}
	out := repoTag(repo, tag)
	info, err := dock.InspectImage(env.dock, out)
	if err != nil {
		return nil, errcode.Annotate(err, "inspect image")
	}

	var repoDigests []string
	digestPrefix := srcRepo + "@"
	for _, digest := range info.RepoDigests {
		if strings.HasPrefix(digest, digestPrefix) {
			repoDigests = append(repoDigests, digest)
		}
	}
	if digest != "" {
		digestWant := digestPrefix + digest
		found := false
		for _, digest := range repoDigests {
			if digest == digestWant {
				found = true
				break
			}
		}
		if !found {
			return nil, errcode.Internalf(
				"digest mismatch, got %q, want %q",
				info.RepoDigests, digestWant,
			)
		}
	}

	sum := newDockerSum(repo, tag, info.ID)
	sum.Origin = from
	return sum, nil
}

func (p *dockerPull) build(env *env, opts *buildOpts) error {
	sum, err := p.pull(env)
	if err != nil {
		return err
	}
	out, err := env.prepareOut(p.out)
	if err != nil {
		return errcode.Annotate(err, "prepare sum output")
	}
	if err := jsonutil.WriteFile(out, sum); err != nil {
		return errcode.Annotate(err, "write image sum")
	}

	if p.tarOut != "" {
		log.Printf("Saving %s", p.tarOut)
		out, err := env.prepareOut(p.tarOut)
		if err != nil {
			return errcode.Annotate(err, "prepare tar output")
		}
		if err := dock.SaveImageGz(env.dock, sum.ID, out); err != nil {
			return errcode.Annotate(err, "save image as tar")
		}
	}
	return nil
}

func (p *dockerPull) meta(env *env) (*buildRuleMeta, error) {
	digest, err := makeDigest(ruleDockerPull, p.name, p.rule)
	if err != nil {
		return nil, errcode.Annotate(err, "digest")
	}

	outs := []string{p.out}
	if p.tarOut != "" {
		outs = append(outs, p.tarOut)
	}

	return &buildRuleMeta{
		name:      p.name,
		outs:      outs,
		dockerOut: true,
		digest:    digest,
	}, nil
}
