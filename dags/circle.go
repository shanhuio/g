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

type searchNode struct {
	start  *MapNode
	this   *MapNode
	last   *searchNode
	length int
}

func traceCircle(trace []*searchNode, snode *searchNode) []*MapNode {
	n := snode.length
	ret := make([]*MapNode, n)
	for i := 0; i < n; i++ {
		ret[n-1-i] = snode.this
		snode = snode.last
	}

	if snode != nil {
		panic("bug")
	}
	return ret
}

func minCircle(nodes map[string]*MapNode) []*MapNode {
	/*
		Breath first search starting from all nodes, where the first node must
		be the smallest node of the circle. Disregarding the complexity of the
		hash maps, Worst case complexity of this algorithm is no larger than
		O(dn^2) where d is the number of edges per node. In practice, a
		dependency map often has a constant number of (d), i.e. importing
		limited number of packages, and hence the complexity is acceptable.
		Also it returns when ths smallest circle is found, so the actual run
		time would often be much smaller than O(dn^2).
	*/

	var trace []*searchNode
	visited := make(map[string]map[string]bool)
	for _, node := range nodes {
		m := make(map[string]bool)
		m[node.Name] = true
		visited[node.Name] = m
	}

	for _, node := range nodes {
		trace = append(trace, &searchNode{
			start:  node,
			this:   node,
			last:   nil,
			length: 1,
		})
	}

	pt := 0
	for pt < len(trace) {
		snode := trace[pt]
		start := snode.start
		vmap := visited[start.Name]
		for name, out := range snode.this.Outs {
			if name == start.Name {
				return traceCircle(trace, snode)
			}

			if vmap[name] {
				// visited before from this start
				continue
			}
			if name < start.Name {
				// a node smaller than the start; skip it
				continue
			}

			trace = append(trace, &searchNode{
				start:  start,
				this:   out,
				last:   snode,
				length: snode.length + 1,
			})
		}

		pt++ // next one
	}

	return nil // no circle found
}
