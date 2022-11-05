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

// CheckDAG checks if a graph is a valid DAG.  It returns true
// when all the graph, links are valid and has no circular
// dependency.
func CheckDAG(g *Graph) error {
	m, err := initMap(g)
	if err != nil {
		return err
	}

	_, err = m.makeLayers()
	return err
}
