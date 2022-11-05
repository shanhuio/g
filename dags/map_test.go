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
	"testing"

	"fmt"
)

func makeFullTestGraph(n int) (*Graph, []string) {
	names := make([]string, n)
	for i := 0; i < n; i++ {
		names[i] = fmt.Sprintf("node-%04d", i)
	}

	nodes := make(map[string][]string)
	for _, name := range names {
		nodes[name] = nil
	}
	for i := 0; i < n; i++ {
		from := names[i]
		for j := i + 1; j < n; j++ {
			to := names[j]
			nodes[from] = append(nodes[from], to)
		}
	}

	return NewGraph(nodes), names
}

func checkSameNodes(t *testing.T, want, got map[string]*MapNode) {
	if len(want) != len(got) {
		t.Errorf("want %d nodes, got %d", len(want), len(got))
		return
	}

	for name, node := range want {
		nodeGot, found := got[name]
		if !found {
			t.Errorf("wanted node %q not found", name)
		}
		if nodeGot != node {
			t.Errorf(
				"want node %q at index %q, got %q",
				node.Name, name, nodeGot.Name,
			)
		}
	}
}

func checkEmptyNodes(t *testing.T, got map[string]*MapNode) {
	if len(got) != 0 {
		t.Errorf("want no nodes, got %d", len(got))
	}
}

func TestMap(t *testing.T) {
	const n = 10
	g, names := makeFullTestGraph(n)
	m, err := NewMap(g)
	if err != nil {
		t.Fatal(err)
	}

	if want := n * (n - 1) / 2; m.Nedge != want {
		t.Errorf("want edge count %d, got %d", want, m.Nedge)
	}
	if want := n - 1; m.Ncrit != want {
		t.Errorf("want critical edge count %d, got %d", want, m.Ncrit)
	}
	if m.Nlayer != n {
		t.Errorf("want layer count %d, got %d", n, m.Nlayer)
	}

	before := make(map[string]*MapNode)
	after := make(map[string]*MapNode)
	for name, node := range m.Nodes {
		if node.Name != name {
			t.Errorf("node at index %q has name %q", name, node.Name)
		}
		after[node.Name] = node
	}

	sortedNodes := m.SortedNodes()
	for i, node := range sortedNodes {
		name := node.Name
		if name != names[i] {
			t.Errorf("node %d, name want %q, got %q", i, names[i], name)
		}
		delete(after, name)
		checkSameNodes(t, before, node.Ins)
		checkSameNodes(t, after, node.Outs)
		checkSameNodes(t, before, node.AllIns)
		checkSameNodes(t, after, node.AllOuts)

		if i == 0 {
			checkEmptyNodes(t, node.CritIns)
		} else {
			last := m.Nodes[names[i-1]]
			want := map[string]*MapNode{last.Name: last}
			checkSameNodes(t, want, node.CritIns)
		}

		if i == n-1 {
			checkEmptyNodes(t, node.CritOuts)
		} else {
			next := m.Nodes[names[i+1]]
			want := map[string]*MapNode{next.Name: next}
			checkSameNodes(t, want, node.CritOuts)
		}

		before[name] = node
	}

	layers := m.SortedLayers()
	if len(layers) != n {
		t.Errorf("want %d layers, got %d", n, len(layers))
	}
	for i, layer := range layers {
		if len(layer) != 1 {
			t.Errorf("want one node in layer, got %d", len(layer))
		}
		if layer[0] != sortedNodes[i] {
			t.Errorf("node in layer %d is incorrect", i)
		}
	}
}
