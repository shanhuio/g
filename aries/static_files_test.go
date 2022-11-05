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

	"net/http/httptest"

	"shanhu.io/pub/httputil"
)

func TestStaticFiles(t *testing.T) {
	static := NewStaticFiles("testdata/static")

	s := httptest.NewServer(Serve(static))
	defer s.Close()

	c := httputil.NewClientMust(s.URL)
	for _, test := range []struct {
		p, want string
	}{
		{"/f1.html", "hello\n"},
		{"/f2.html", "hi\n"},
	} {
		reply, err := c.GetString(test.p)
		if err != nil {
			t.Errorf("%q - got error: %s", test.p, err)
			continue
		}
		if reply != test.want {
			t.Errorf("%q - want %q, got %q", test.p, test.want, reply)
		}
	}
}
