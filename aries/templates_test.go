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

func TestTemplates(t *testing.T) {
	tmpls := NewTemplates("testdata/templates", nil)

	f := func(c *C) error {
		dat := struct {
			Message1, Message2 string
		}{
			Message1: "hello",
			Message2: "hi",
		}
		return tmpls.Serve(c, c.Rel(), dat)
	}

	s := httptest.NewServer(Func(f))
	defer s.Close()

	c := httputil.NewClientMust(s.URL)
	for _, test := range []struct {
		url, want string
	}{
		{"/t1.html", "hello\n"},
		{"/t2.html", "hi\n"},
	} {
		reply, err := c.GetString(test.url)
		if err != nil {
			t.Errorf(
				"http get %q, got error: %s",
				test.url, err,
			)
		}
		if reply != test.want {
			t.Errorf(
				"http get %q, want %q, got %q",
				test.url, test.want, reply,
			)
		}
	}
}

func TestTemplatesJSON(t *testing.T) {
	tmpls := NewTemplates(TemplatesJSON, nil)

	f := func(c *C) error {
		dat := struct {
			Message1, Message2 string
		}{
			Message1: "hello",
			Message2: "hi",
		}
		return tmpls.Serve(c, c.Rel(), dat)
	}

	s := httptest.NewServer(Func(f))
	defer s.Close()

	c := httputil.NewClientMust(s.URL)
	test := struct {
		url, want string
	}{
		"/", `{"Message1":"hello","Message2":"hi"}`,
	}

	reply, err := c.GetString(test.url)
	if err != nil {
		t.Errorf(
			"http get %q, got error: %s",
			test.url, err,
		)
	}
	if reply != test.want {
		t.Errorf(
			"http get %q, want %q, got %q",
			test.url, test.want, reply,
		)
	}
}
