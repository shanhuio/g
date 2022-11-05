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

	"fmt"
	"net/http/httptest"

	"shanhu.io/pub/httputil"
)

func makeEchoRel(s string) Func {
	return func(c *C) error {
		fmt.Fprintf(c.Resp, "%s: %s", s, c.Rel())
		return nil
	}
}

func TestRouter(t *testing.T) {
	r := NewRouter()
	r.Get("something", StringFunc("xxx"))
	r.Dir("books", makeEchoRel("books"))

	s := httptest.NewServer(Serve(r))
	defer s.Close()

	c := s.Client()
	host := s.URL

	for _, test := range []struct {
		p, want string
	}{
		{"/something", "xxx"},
		{"/books/xxx", "books: xxx"},
		{"/books/yyy", "books: yyy"},
		{"/books/", "books: "},
		{"/books", "books: "},
	} {
		got, err := httputil.GetString(c, host+test.p)
		if err != nil {
			t.Errorf("get %q, got error: %s", test.p, err)
			continue
		}

		if got != test.want {
			t.Errorf(
				"get %q, want %q in response, got %q",
				test.p, test.want, got,
			)
		}
	}

	for _, p := range []string{
		"/something/xxx",
		"/bookss",
		"/something/",
	} {
		code, err := httputil.GetCode(c, host+p)
		if err != nil {
			t.Error(err)
			continue
		}

		if code != 404 {
			t.Errorf("get %q, want 404 response, got %d", p, code)
		}
	}
}

func TestRouterWithIndex(t *testing.T) {
	r := NewRouter()
	r.Index(StringFunc("index"))

	sub := NewRouter()
	sub.Index(StringFunc("sub-index"))
	r.Dir("sub", sub.Serve)

	s := httptest.NewServer(Serve(r))
	defer s.Close()

	for _, test := range []struct {
		p, want string
	}{
		{"", "index"},
		{"/", "index"},
		{"/sub", "sub-index"},
		{"/sub/", "sub-index"},
	} {
		got, err := httputil.GetString(s.Client(), s.URL+test.p)
		if err != nil {
			t.Error(err)
			return
		}
		if got != test.want {
			t.Errorf(
				"get index page, want %q in response, got %q",
				test.want, got,
			)
		}
	}
}

func TestRouterWithDefault(t *testing.T) {
	r := NewRouter()
	r.Index(StringFunc("index"))
	r.Default(StringFunc("default"))
	s := httptest.NewServer(Serve(r))
	defer s.Close()

	got, err := httputil.GetString(s.Client(), s.URL+"/notfound")
	if err != nil {
		t.Error(err)
		return
	}

	if got != "default" {
		t.Errorf(
			"get a 404 page, want %q in response, got %q",
			"default", got,
		)
	}
}

func TestRouterWithRedirect(t *testing.T) {
	r := NewRouter()
	r.Index(StringFunc("index"))
	r.Get("redirect", RedirectTo("/"))
	s := httptest.NewServer(Serve(r))
	defer s.Close()

	got, err := httputil.GetString(s.Client(), s.URL+"/redirect")
	if err != nil {
		t.Error(err)
		return
	}

	if got != "index" {
		t.Errorf("got %q, want index", got)
	}
}
