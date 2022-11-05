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

// MapNodeView is the position of a map node.
type MapNodeView struct {
	Name        string
	DisplayName string
	X, Y        int
	CritIns     []string
	CritOuts    []string
}

// MapView is a layout view of a Map.
type MapView struct {
	Height    int
	Width     int
	Nodes     map[string]*MapNodeView
	IsTopDown bool
}

// AssignDisplayName assigns display names in a map.
func (v *MapView) AssignDisplayName(f func(s string) string) {
	for _, n := range v.Nodes {
		n.DisplayName = f(n.Name)
	}
}

// Reverse reverse a map view.
func (v *MapView) Reverse() {
	for _, node := range v.Nodes {
		node.X = v.Width - 1 - node.X
		node.CritIns, node.CritOuts = node.CritOuts, node.CritIns
	}

	v.IsTopDown = !v.IsTopDown
}
