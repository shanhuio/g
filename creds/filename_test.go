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

package creds

import (
	"testing"
)

func TestFilename(t *testing.T) {
	o := func(from, to string) {
		got := Filename(from)
		if got != to {
			t.Errorf("Filename(%q) mapped to %q, want %q", from, got, to)
		}
	}

	o("shanhu.io", "shanhu-io")
	o("smallrepo.com", "smallrepo-com")
	o("localhost:3356", "localhost-3356")
	o("localhost:3335", "localhost-3335")
}
