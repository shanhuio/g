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
	"shanhu.io/g/lexing"
)

const (
	nodeSrc  = "src"
	nodeRule = "rule"
	nodeOut  = "out"
	nodeRun  = "run"
	nodeSub  = "sub"
)

type buildNode struct {
	name string
	typ  string
	deps []string
	pos  *lexing.Pos

	ruleType string
	rule     buildRule
	ruleMeta *buildRuleMeta

	sub *subBuilds
}

func (n *buildNode) mainOut() string {
	if m := n.ruleMeta; m != nil && len(m.outs) > 0 {
		return m.outs[0]
	}
	return ""
}
