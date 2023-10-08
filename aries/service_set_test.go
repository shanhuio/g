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

package aries

import (
	"testing"

	"fmt"
	"net/http/httptest"

	"shanhu.io/g/httputil"
)

type testAuth struct{}

func (a *testAuth) Serve(c *C) error {
	return Miss
}

func (a *testAuth) Setup(c *C) error {
	bearer := Bearer(c)
	if bearer == "test-user" {
		c.User = "test-user"
	} else if bearer == "test-admin" {
		c.User = "test-admin"
		c.UserLevel = 1
	}
	return nil
}

func TestServiceSet(t *testing.T) {
	set := &ServiceSet{
		Auth: new(testAuth),
		Resource: Func(func(c *C) error {
			if c.Path == "/style" {
				fmt.Fprint(c.Resp, "style")
				return nil
			}
			return Miss
		}),
		Guest: Func(func(c *C) error {
			if c.Path == "/guest" {
				fmt.Fprint(c.Resp, "guest")
				return nil
			} else if c.Path == "/username" {
				fmt.Fprint(c.Resp, c.User)
				return nil
			}
			return Miss
		}),
		User: Func(func(c *C) error {
			if c.Path == "/" {
				fmt.Fprint(c.Resp, "user")
				return nil
			}
			return Miss
		}),
		Admin: Func(func(c *C) error {
			if c.Path == "/admin" {
				fmt.Fprint(c.Resp, "admin")
				return nil
			}
			return Miss
		}),
	}

	s := httptest.NewServer(Serve(set))
	defer s.Close()

	type testCase struct {
		p, token string
		wantCode int
		want     string
	}

	runTest := func(url string, test *testCase) {
		c := httputil.NewTokenClientMust(url, test.token)
		if test.wantCode != 200 {
			got, err := c.GetCode(test.p)
			if err != nil {
				t.Errorf("%q@%s - got error: %s", test.p, test.token, err)
				return
			}
			if got != test.wantCode {
				t.Errorf(
					"%q@%s - want %d, got %d",
					test.p, test.token, test.wantCode, got,
				)
			}
			return
		}

		got, err := c.GetString(test.p)
		if err != nil {
			t.Errorf("%q@%s - got error: %s", test.p, test.token, err)
			return
		}
		if got != test.want {
			t.Errorf(
				"%q@%s - want %q, got %q",
				test.p, test.token, test.want, got,
			)
		}
	}

	for _, test := range []*testCase{
		{"/style", "", 200, "style"},
		{"/username", "", 200, ""},
		{"/username", "test-user", 200, "test-user"},
		{"/username", "test-admin", 200, "test-admin"},
		{"/guest", "", 200, "guest"},
		{"/guest", "test-user", 200, "guest"},
		{"/", "test-user", 200, "user"},
		{"/", "test-admin", 200, "user"},
		{"/admin", "test-user", 404, ""},
		{"/admin", "test-admin", 200, "admin"},
		{"/not-found", "", 404, ""},
	} {
		runTest(s.URL, test)
	}

	si := httptest.NewServer(Serve(Func(set.ServeInternal)))
	defer si.Close()

	for _, test := range []*testCase{
		{"/style", "", 200, "style"},
		{"/username", "", 403, ""},
		{"/username", "test-user", 403, "test-user"},
		{"/username", "test-admin", 200, "test-admin"},
		{"/guest", "", 403, ""},
		{"/guest", "test-user", 403, "guest"},
		{"/", "test-user", 403, "user"},
		{"/", "test-admin", 200, "user"},
		{"/admin", "test-user", 403, ""},
		{"/admin", "test-admin", 200, "admin"},
		{"/not-found", "", 403, ""},
	} {
		runTest(si.URL, test)
	}
}
