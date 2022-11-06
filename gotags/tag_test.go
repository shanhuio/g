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

package gotags

import (
	"testing"
)

func TestTagString(t *testing.T) {
	tag := NewTag("tagname", "filename", 2, "x")
	tag.Fields["access"] = "public"
	tag.Fields["type"] = "struct"
	tag.Fields["signature"] = "()"
	tag.Fields["empty"] = ""

	expected := "tagname\tfilename\t2;\"\tx\taccess:public\tline:2\tsignature:()\ttype:struct"

	s := tag.String()
	if s != expected {
		t.Errorf("Tag.String()\n  is:%s\nwant:%s", s, expected)
	}
}
