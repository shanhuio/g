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

// MapNode is a node in the DAG graph
type MapNode struct {
	Name string

	Ins  map[string]*MapNode // direct input nodes
	Outs map[string]*MapNode // direct output nodes

	AllIns  map[string]*MapNode // direct and indirect input nodes
	AllOuts map[string]*MapNode // direct and indirect output nodes

	// critical nodes is the minimum set of nodes that keeps
	// the same partial order of the nodes in the DAG graph
	CritIns  map[string]*MapNode // critical input nodes
	CritOuts map[string]*MapNode // critical output nodes

	x, y     int // for calculating layout position
	layer    int // min layer
	newLayer int // new layer after pushing

	nhit int // for counting on layers
}

func newMapNode(name string) *MapNode {
	return &MapNode{
		Name:     name,
		AllIns:   make(map[string]*MapNode),
		AllOuts:  make(map[string]*MapNode),
		CritIns:  make(map[string]*MapNode),
		CritOuts: make(map[string]*MapNode),
		Ins:      make(map[string]*MapNode),
		Outs:     make(map[string]*MapNode),
	}
}
