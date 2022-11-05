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
	"math/rand"
)

func makeRandomDAG(n int, p float64) *Graph {
	names := make([]string, n)
	for i := 0; i < n; i++ {
		names[i] = fmt.Sprintf("node%d", i)
	}

	nodes := make(map[string][]string)
	for i := 0; i < n; i++ {
		from := names[i]
		var edges []string
		for j := i + 1; j < n; j++ {
			to := names[j]
			if rand.Float64() < p {
				edges = append(edges, to)
			}
		}
		nodes[from] = edges
	}

	return NewGraph(nodes)
}

func TestCheckDAG(t *testing.T) {
	for _, test := range []struct {
		n int
		p float64
	}{
		{10, 0.2},
		{10, 0},
		{10, 1.0},
		{50, 0.2},
		{50, 0.5},
		{50, 1.0},
	} {
		g := makeRandomDAG(test.n, test.p)
		if err := CheckDAG(g); err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	}
}

func TestCheckDAGCircle(t *testing.T) {
	for _, n := range []int{2, 5, 10, 50} {
		names := make([]string, n)
		for i := 0; i < n; i++ {
			names[i] = fmt.Sprintf("node%d", i)
		}

		nodes := make(map[string][]string)
		for i := 0; i < n; i++ {
			if i == 0 {
				nodes[names[i]] = []string{names[n-1]}
			} else {
				nodes[names[i]] = []string{names[i-1]}
			}
		}

		g := NewGraph(nodes)
		err := CheckDAG(g)
		if err == nil {
			t.Error("a circle is not a DAG")
		}
	}
}
