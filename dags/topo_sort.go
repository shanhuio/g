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

// TopoSort topologically sorts the graph.
func TopoSort(g *Graph) ([]string, error) {
	m, err := NewMap(g)
	if err != nil {
		return nil, err
	}

	nodes := m.SortedNodes()
	ret := make([]string, 0, len(nodes))
	for _, node := range nodes {
		ret = append(ret, node.Name)
	}

	return ret, nil
}
