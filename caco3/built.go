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
	"shanhu.io/pub/dock"
	"shanhu.io/pub/errcode"
)

type built struct {
	Outs    []*fileStat  `json:",omitempty"` // A list of outputs.
	Dockers []*dockerSum `json:",omitempty"` // A contaienr image.
}

func newBuilt(env *env, meta *buildRuleMeta) (*built, error) {
	b := new(built)
	for i, out := range meta.outs {
		if i == 0 && meta.dockerOut {
			sum, err := loadDockerSum(env.out(out))
			if err != nil {
				return nil, errcode.Annotatef(
					err, "read docker sum: %s", out,
				)
			}
			b.Dockers = append(b.Dockers, sum)
		}
		stat, err := newOutFileStat(env, out)
		if err != nil {
			return nil, errcode.Annotatef(
				err, "get output stat: %s", out,
			)
		}
		b.Outs = append(b.Outs, stat)
	}
	return b, nil
}

func checkSameBuilt(env *env, b *built) (bool, error) {
	for _, out := range b.Outs {
		same, err := sameFileStat(env, out)
		if err != nil {
			return false, errcode.Annotatef(
				err, "check output stat of %q", out.Name,
			)
		}
		if !same {
			return false, nil
		}
	}

	for _, d := range b.Dockers {
		repoTag := repoTag(d.Repo, d.Tag)
		info, err := dock.InspectImage(env.dock, repoTag)
		if err != nil {
			if errcode.IsNotFound(err) {
				return false, nil // Image not found.
			}
			return false, errcode.Annotatef(err, "inspect docker %s", repoTag)
		}
		if info.ID != d.ID {
			return false, nil // Image ID changed.
		}
	}

	return true, nil
}
