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

package objects

import (
	"testing"
)

func TestJSON(t *testing.T) {
	b := NewMem()

	type data struct {
		S string
	}
	d := &data{"something"}

	hash, err := CreateJSON(b, d)
	if err != nil {
		t.Fatal(err)
	}

	got := new(data)
	if err := ReadJSON(b, hash, got); err != nil {
		t.Fatal(err)
	}
	if got.S != d.S {
		t.Errorf("got %q, want %q", got.S, d.S)
	}
}
