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

package lexing

import (
	"testing"
)

func TestIsPkgName(t *testing.T) {
	testName := []string{"a", "abc", "a12", "a12b"}
	for _, name := range testName {
		if !IsPkgName(name) {
			t.Errorf("%v should be a package name", name)
		}
	}
	testName = []string{"aB", "1abc", "", "  ", "a~", "%a", "A1", "TBC", "$abc1"}
	for _, name := range testName {
		if IsPkgName(name) {
			t.Errorf("%v should not be a package name", name)
		}
	}
}
