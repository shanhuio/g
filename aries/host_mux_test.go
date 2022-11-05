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

	"shanhu.io/aries/https/httpstest"
	"shanhu.io/pub/httputil"
)

func TestHostMux(t *testing.T) {
	m := NewHostMux()
	m.Set("shanhu.io", StringFunc("shanhu"))
	m.Set("h8liu.io", StringFunc("h8liu"))

	s, err := httpstest.NewServer([]string{
		"shanhu.io", "h8liu.io",
	}, Serve(m))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	c := s.Client()

	for _, test := range []struct {
		url, want string
	}{
		{"https://shanhu.io", "shanhu"},
		{"https://h8liu.io", "h8liu"},
	} {
		got, err := httputil.GetString(c, test.url)
		if err != nil {
			t.Errorf("get %q, got error %s", test.url, err)
		} else if got != test.want {
			t.Errorf("get %q, got %q, want %q", test.url, got, test.want)
		}
	}
}
