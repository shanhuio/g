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

	"reflect"
	"strings"
)

func checkGraphEqual(t *testing.T, got *Graph, nodes map[string][]string) {
	want := NewGraph(nodes)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestGraphReverse(t *testing.T) {
	g := NewGraph(map[string][]string{
		"a": {"b", "c"},
		"b": {"c"},
		"c": nil,
	})

	got := g.Reverse()
	checkGraphEqual(t, got, map[string][]string{
		"c": {"a", "b"},
		"b": {"a"},
		"a": nil,
	})
}

func TestGraphRemove(t *testing.T) {
	g := NewGraph(map[string][]string{
		"a": {"b", "c"},
		"b": {"c"},
		"c": nil,
	})

	got := g.Remove("b")
	checkGraphEqual(t, got, map[string][]string{
		"a": {"c"},
		"c": nil,
	})
}

func TestSubGraph(t *testing.T) {
	g := NewGraph(map[string][]string{
		"a": {"b", "c"},
		"b": {"c"},
		"c": nil,
	})

	got := g.SubGraph(func(n string) bool {
		return n != "b"
	})

	checkGraphEqual(t, got, map[string][]string{
		"a": {"c"},
		"c": nil,
	})
}

func TestGraphRename(t *testing.T) {
	g := NewGraph(map[string][]string{
		"a": {"b", "c"},
		"b": {"c"},
		"c": nil,
	})

	got, err := g.Rename(func(name string) (string, error) {
		return strings.ToUpper(name), nil
	})
	if err != nil {
		t.Fatal(err)
	}

	checkGraphEqual(t, got, map[string][]string{
		"A": {"B", "C"},
		"B": {"C"},
		"C": nil,
	})
}
