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

type bundle struct {
	rule *Bundle
	name string
	deps []string
}

func newBundle(env *env, p string, r *Bundle) *bundle {
	name := makeRelPath(p, r.Name)
	var deps []string
	for _, dep := range r.Deps {
		deps = append(deps, makePath(p, dep))
	}

	return &bundle{
		name: name,
		deps: deps,
		rule: r,
	}
}

func (b *bundle) build(env *env, opts *buildOpts) error {
	return nil
}

func (b *bundle) meta(env *env) (*buildRuleMeta, error) {
	d, err := makeDigest(ruleBundle, b.name, struct{}{})
	if err != nil {
		return nil, errcode.Annotate(err, "digest")
	}

	return &buildRuleMeta{
		name:   b.name,
		deps:   b.deps,
		digest: d,
	}, nil
}
