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

	"net/http/httptest"

	"shanhu.io/pub/errcode"
	"shanhu.io/pub/httputil"
)

func TestJSONString(t *testing.T) {
	const msg = "hello"
	const reply = "hi"

	f := func(c *C) error {
		var s string
		if err := UnmarshalJSONBody(c, &s); err != nil {
			return err
		}
		if s != msg {
			return errcode.InvalidArgf("not the right message")
		}
		return ReplyJSON(c, reply)
	}

	s := httptest.NewServer(Func(f))
	defer s.Close()

	c := httputil.NewClientMust(s.URL)
	var str string
	if err := c.JSONCall("/", msg, &str); err != nil {
		t.Fatal(err)
	}

	if str != reply {
		t.Errorf("want %q, got %q", reply, str)
	}
}

func TestJSONStruct(t *testing.T) {
	type data struct {
		Message string
	}
	const msg = "hello"
	const reply = "hi"

	f := func(c *C) error {
		d := new(data)
		if err := UnmarshalJSONBody(c, d); err != nil {
			return err
		}
		if d.Message != msg {
			return errcode.InvalidArgf("not the right message")
		}
		return ReplyJSON(c, &data{Message: reply})
	}

	s := httptest.NewServer(Func(f))
	defer s.Close()

	c := httputil.NewClientMust(s.URL)
	d := new(data)
	if err := c.JSONCall("/", &data{Message: msg}, d); err != nil {
		t.Fatal(err)
	}

	if d.Message != reply {
		t.Errorf("want %q, got %q", reply, d.Message)
	}
}
