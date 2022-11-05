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
	"shanhu.io/pub/errcode"
)

func buildNodeDigest(
	env *env, n *buildNode, deps map[string]string,
) (string, error) {
	switch n.typ {
	case nodeRule:
		action := &buildAction{
			Deps:     deps,
			RuleType: n.ruleType,
		}
		if meta := n.ruleMeta; meta != nil {
			if meta.digest == "" {
				return "", nil
			}
			action.Rule = meta.digest
			action.Outs = meta.outs
			action.DockerOut = meta.dockerOut
		}
		d, err := makeDigest("build_action", "", action)
		if err != nil {
			return "", errcode.Annotate(err, "digest build action")
		}
		return d, nil
	case nodeSrc:
		stat, err := newSrcFileStat(env, n.name)
		if err != nil {
			return "", errcode.Annotatef(err, "stat file %q", n.name)
		}
		d, err := makeDigest("src", "", stat)
		if err != nil {
			return "", errcode.Annotate(err, "digest source file")
		}
		return d, nil
	case nodeOut:
		action := &buildAction{
			Deps:     deps,
			OutputOf: n.name,
		}
		d, err := makeDigest("out", "", action)
		if err != nil {
			return "", errcode.Annotate(err, "digest output-of")
		}
		return d, nil
	default:
		return "", nil
	}
}
