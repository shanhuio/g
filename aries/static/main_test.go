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

package static

import (
	"testing"

	"shanhu.io/g/aries/ariestest"
	"shanhu.io/g/httputil"
)

func TestMain(t *testing.T) {
	service := makeService("testdata")
	s, err := ariestest.HTTPSServer("shanhu.io", service)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	c := httputil.NewClientMust("https://shanhu.io")
	c.Transport = s.Transport

	str, err := c.GetString("/")
	if err != nil {
		t.Fatal(err)
	}
	const want = "hello\n"
	if str != want {
		t.Errorf("get /, want %q, got %q", want, str)
	}
}
