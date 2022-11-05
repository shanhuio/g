// Copyright (C) 2022  Shanhu Tech Inc.
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

package aries

import (
	"testing"
)

func TestTrie(t *testing.T) {
	root := newTrieRoot()
	add := func(input string, added bool) {
		if added != root.add(input) {
			t.Errorf("expected dulplicated add = %v, got %v", added, !added)
		}
	}
	find := func(input, pref string, match bool) {
		rp, rm := trieFind(root, input)
		if rp != pref {
			t.Errorf("expected pref = %q, got %q", pref, rp)
		}
		if rm != match {
			t.Errorf("expected match = %v, got %v", match, rm)
		}
	}
	add("", false)
	find("a", "", false)
	add("axy45678", true)
	add("abc", true)
	add("axy", true)
	find("abc", "abc", true)
	find("abcd", "abc", false)
	find("a", "", false)
	find("ax", "", false)
	add("abc", false)
	add("a", true)
	find("a", "a", true)
	find("axy12", "axy", false)
	add("axy120", true)
	find("axy12", "axy", false)
	find("", "", true)
	add("ax456", true)
	add("ax4567", true)
	find("ax45", "a", false)
	find("dtc", "", false)
	find("axy45678", "axy45678", true)
}
