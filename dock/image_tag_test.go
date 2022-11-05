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

package dock

import (
	"testing"
)

func TestImageTag(t *testing.T) {
	for _, test := range []struct {
		s, img, tag string
	}{
		{s: "nextcloud", img: "nextcloud", tag: ""},
		{s: "nextcloud:19", img: "nextcloud", tag: "19"},
		{s: "shanhu.io/doorway", img: "shanhu.io/doorway", tag: ""},
		{s: "shanhu.io/doorway:v1", img: "shanhu.io/doorway", tag: "v1"},
	} {
		img, tag := ParseImageTag(test.s)
		if img != test.img || tag != test.tag {
			t.Errorf(
				"ParseImageTag(%q), got (%q, %q), want (%q, %q)",
				test.s, img, tag, test.img, test.tag,
			)
		}
	}
}
