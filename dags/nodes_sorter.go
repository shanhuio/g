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

type mapNodes []*MapNode

func (s mapNodes) Len() int      { return len(s) }
func (s mapNodes) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type byLayer struct{ mapNodes }

func (s byLayer) Less(i, j int) bool {
	a := s.mapNodes[i]
	b := s.mapNodes[j]

	if a.layer < b.layer {
		return true
	}
	if a.layer > b.layer {
		return false
	}
	return a.Name < b.Name
}

type byNcritOuts struct{ mapNodes }

func (s byNcritOuts) Less(i, j int) bool {
	n1 := s.mapNodes[i]
	n2 := s.mapNodes[j]

	nin1 := len(n1.CritIns)
	nin2 := len(n2.CritIns)
	if nin1 < nin2 {
		return false
	}
	if nin1 > nin2 {
		return true
	}

	nout1 := len(n1.CritOuts)
	nout2 := len(n2.CritOuts)

	if nout1 < nout2 {
		return false
	}
	if nout1 > nout2 {
		return true
	}

	return n1.Name < n2.Name
}

type byName struct{ mapNodes }

func (s byName) Less(i, j int) bool {
	return s.mapNodes[i].Name < s.mapNodes[j].Name
}
