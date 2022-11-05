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

// checkPush returns if a node is able and worthy to upper layer
func checkPush(m *Map, node *MapNode) (able, worthy bool) {
	if node.layer == m.Nlayer-1 {
		return false, false // already at the top layer
	}

	for _, out := range node.CritOuts {
		if out.layer > node.layer+1 {
			// exists an edge that can be shorten after
			// this node is pushed to upper layer
			worthy = true
			continue
		}

		subAble, subWorthy := checkPush(m, out)
		if !subAble {
			return false, false
		}
		if subWorthy {
			worthy = true
		}
	}

	return true, worthy
}

// pushWorthy checks if the node is worthy to push
// a node is push worthy if pushing it can reduce the total
// length of the edges
func pushWorthy(m *Map, node *MapNode) bool {
	_, ret := checkPush(m, node)
	return ret
}

// pushNode pushes a node to upper layer
func pushNode(m *Map, node *MapNode, pushed map[string]*MapNode) {
	// pushing all nodes on the right
	for _, out := range node.CritOuts {
		if out.layer > node.layer+1 {
			continue // no need to push this
		}
		pushNode(m, out, pushed)
	}

	pushed[node.Name] = node // add to the set
}

// pushTight pushes map nodes to upper layers if possible
// so that the nodes are closer to its dependencies
func pushTight(m *Map) {
	nodes := m.SortedNodes()

	n := len(nodes)
	for i := range nodes {
		node := nodes[n-1-i] // in reversed order
		// we push upper layer nodes first

		for pushWorthy(m, node) {
			pushed := make(map[string]*MapNode)
			pushNode(m, node, pushed)

			// update the layer
			for _, p := range pushed {
				p.layer++
				if p.layer >= m.Nlayer {
					panic("pushing to hard, increasing layers")
				}
			}
		}
	}
}
