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

package osutil

import (
	"testing"

	"os"
	"path"
)

func TestIsDir(t *testing.T) {
	d := t.TempDir()

	ok, err := IsDir(d)
	if err != nil {
		t.Fatal("check is dir: ", err)
	}
	if !ok {
		t.Errorf("IsDir(%q) should return true", d)
	}

	f := path.Join(d, "post")
	if err := os.WriteFile(f, []byte("post"), 0600); err != nil {
		t.Fatal("write file: ", err)
	}

	ok, err = IsDir(f)
	if err != nil {
		t.Fatal("check file is dir: ", err)
	}
	if ok {
		t.Errorf("IsDir(%q) should return false", f)
	}

	ok, err = IsDir(path.Join(d, "ghost"))
	if err != nil {
		t.Fatal("check ghost dir: ", err)
	}
	if ok {
		t.Errorf("IsDir(%q) should return false", f)
	}
}
