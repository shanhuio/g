// Copyright (C) 2021  Shanhu Tech Inc.
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

package dags

import (
	"fmt"
)

// Closure makes the minimum sub map of a map that contains all the given
// nodes, and nodes that are in between these nodes.
func Closure(m *Map, nodes []string) *Map {
	nodeSet := make(map[string]bool)
	ins := make(map[string]bool)
	outs := make(map[string]bool)

	for _, name := range nodes {
		node := m.Nodes[name]
		if node == nil {
			panic(fmt.Errorf("%q node not found", name))
		}
		nodeSet[name] = true

		for in := range node.AllIns {
			ins[in] = true
		}
		for out := range node.AllOuts {
			outs[out] = true
		}
	}

	for in := range ins {
		if outs[in] {
			nodeSet[in] = true
		}
	}

	g := make(map[string][]string)
	for name := range nodeSet {
		var outs []string
		node := m.Nodes[name]
		for out := range node.Outs {
			if nodeSet[out] {
				outs = append(outs, out)
			}
		}
		g[name] = outs
	}

	m, err := NewMap(NewGraph(g))
	if err != nil {
		panic(err)
	}
	return m
}
