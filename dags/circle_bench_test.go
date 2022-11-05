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
	"math/rand"
	"strconv"
	"testing"
)

func BenchmarkFindCircle(b *testing.B) {
	for n := 0; n < b.N; n++ {
		benchmarkFindCircle()
	}
}

func benchmarkFindCircle() []*MapNode {
	ret := make(map[string][]string)

	for i := 0; i < 1000; i++ {
		var edge []string
		for j := 0; j < 1000; j++ {
			if i == j {
				continue
			}
			q := rand.Int31n(500)
			if q < 1 {
				var flag bool
				for _, ele := range ret[strconv.Itoa(j)] {
					if ele == strconv.Itoa(i) {
						flag = true
					}
				}
				if !flag {
					edge = append(edge, strconv.Itoa(j))
				}
			}
		}
		ret[strconv.Itoa(i)] = edge
	}

	g := &Graph{Nodes: ret}
	nodes, _ := initMap(g)

	return minCircle(nodes.Nodes)
}
