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

package srpc

import (
	"net/url"
	"testing"
)

func TestPathJoin(t *testing.T) {
	for _, test := range []struct {
		server, p, want string
	}{
		{server: "http://h", p: "p", want: "http://h/p"},
		{server: "http://h/", p: "p", want: "http://h/p"},
		{server: "http://h", p: "/p", want: "http://h/p"},
		{server: "http://h/", p: "/p", want: "http://h/p"},
		{server: "http://h/p", p: "q", want: "http://h/p/q"},
		{server: "http://h/p", p: "/q", want: "http://h/p/q"},
		{server: "http://h/p/", p: "q", want: "http://h/p/q"},
		{server: "http://h/p/", p: "/q", want: "http://h/p/q"},
		{server: "http://h/p/", p: "q/", want: "http://h/p/q/"},
		{server: "http://h/p", p: "q/", want: "http://h/p/q/"},
	} {
		server, err := url.Parse(test.server)
		if err != nil {
			t.Fatalf("parse server %q: %s", test.server, err)
		}
		got := urlJoin(server, test.p).String()
		if got != test.want {
			t.Errorf(
				"pathJoin(%q, %q) = %q, want %q",
				test.server, test.p, got, test.want,
			)
		}
	}
}
