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

package srpc_test

import (
	"context"
	"net/http/httptest"
	"net/url"
	"testing"

	"shanhu.io/g/aries"
	"shanhu.io/g/srpc"
)

func TestCaller(t *testing.T) {
	r := aries.NewRouter()
	r.Call("echo", func(c *aries.C, req string) (string, error) {
		return req, nil
	})

	s := httptest.NewServer(aries.Serve(r))
	defer s.Close()

	serverURL, err := url.Parse(s.URL)
	if err != nil {
		t.Fatal("parse server url:", err)
	}
	t.Log("server url:", serverURL)

	ctx := context.Background()

	c := srpc.NewCaller(serverURL)

	const msg = "hello"
	var resp string
	if err := c.Call(ctx, "echo", msg, &resp); err != nil {
		t.Fatal("call:", err)
	}
	if resp != msg {
		t.Errorf("got response %q, want %q", resp, msg)
	}
}
