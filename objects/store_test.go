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

	"shanhu.io/pub/errcode"
	"shanhu.io/pub/hashutil"
)

func testStore(t *testing.T, s Store) {
	if _, err := s.Get("not really a hash"); !errcode.IsNotFound(err) {
		t.Errorf("got %s, want not found", err)
	}

	const dat = "some data here"
	bs := []byte(dat)
	h := hashutil.Hash(bs)
	if _, err := s.Get(h); !errcode.IsNotFound(err) {
		t.Errorf("got %s, want not found", err)
	}

	k, err := s.Put(bs)
	if err != nil {
		t.Fatal(err)
	}
	if k != h {
		t.Errorf("got %q, want %q", k, h)
	}

	readBack, err := s.Get(k)
	if err != nil {
		t.Fatal(err)
	}
	got := string(readBack)
	if got != dat {
		t.Errorf("got %q, want %q", got, dat)
	}
}
