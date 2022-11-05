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

func TestRoute(t *testing.T) {
	for _, test := range []struct {
		path    string
		cleaned string
		size    int
		isDir   bool
	}{
		{"/", "", 0, true},
		{"/something", "/something", 1, false},
		{"/something/", "/something", 1, true},
		{"/a/b/c", "/a/b/c", 3, false},
		{"/a//c", "/a/c", 2, false},
		{"/////", "", 0, true},
	} {
		r := newRoute(test.path)
		got := r.path()
		if got != test.cleaned {
			t.Errorf(
				"clean route for %q, want %q, got %q",
				test.path, test.cleaned, got,
			)
		}

		size := r.size()
		if size != test.size {
			t.Errorf(
				"route size for %q, want %d, got %d",
				test.path, test.size, size,
			)
		}

		if r.isDir != test.isDir {
			t.Errorf(
				"route for %q, want isDir=%t, got %t",
				test.path, test.isDir, r.isDir,
			)
		}
	}
}
