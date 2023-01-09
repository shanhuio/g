// Copyright (C) 2023  Shanhu Tech Inc.
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

	"bytes"
	"io"

	"shanhu.io/pub/errcode"
	"shanhu.io/pub/hashutil"
)

func testObjects(t *testing.T, b Objects) {
	if _, err := b.Open("not a hash"); !errcode.IsNotFound(err) {
		t.Errorf("got %s, want not found", err)
	}

	const dat = "some data here"
	bs := []byte(dat)
	h := hashutil.Hash(bs)
	if _, err := b.Open(h); !errcode.IsNotFound(err) {
		t.Errorf("got %s, want not found", err)
	}

	if has, err := b.Has(h); err != nil {
		t.Fatal(err)
	} else if has {
		t.Errorf("Has(%q), got true, want false", h)
	}

	k, err := b.Create(bytes.NewReader(bs))
	if err != nil {
		t.Fatal(err)
	}
	if k != h {
		t.Errorf("got %q, want %q", k, h)
	}

	if has, err := b.Has(h); err != nil {
		t.Fatal(err)
	} else if !has {
		t.Errorf("Has(%q), got false, want true", h)
	}

	r, err := b.Open(k)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	bs, err = io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}

	got := string(bs)
	if got != dat {
		t.Errorf("read back got %q, want %q", got, dat)
	}
}
