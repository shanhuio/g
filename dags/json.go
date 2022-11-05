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
	"encoding/json"
)

// N is a node in the minified DAG visualization result.
type N struct {
	N    string   `json:"n"`
	F    string   `json:"f,omitempty"`
	X    int      `json:"x"`
	Y    int      `json:"y"`
	Ins  []string `json:"i"`
	Outs []string `json:"o"`
}

// M is a node in the minified DAG visualization result.
type M struct {
	Height int           `json:"h"`
	Width  int           `json:"w"`
	Nodes  map[string]*N `json:"n"`
}

// Output returns a json'able object of a map.
func Output(v *MapView) *M {
	res := &M{
		Height: v.Height,
		Width:  v.Width,
		Nodes:  make(map[string]*N),
	}

	for name, node := range v.Nodes {
		display := node.DisplayName
		if display == "" {
			display = name
		}

		n := &N{
			N:    display,
			F:    name,
			X:    node.X,
			Y:    node.Y,
			Ins:  node.CritIns,
			Outs: node.CritOuts,
		}

		res.Nodes[name] = n
	}

	return res
}

func marshalMap(m *MapView) []byte {
	res := Output(m)
	ret, e := json.MarshalIndent(res, "", "    ")
	if e != nil {
		panic(e)
	}

	return ret
}
