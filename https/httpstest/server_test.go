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

package httpstest

import (
	"testing"

	"io/ioutil"
	"net/http"

	"shanhu.io/pub/aries"
)

func checkBody(t *testing.T, resp *http.Response, msg string) {
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %s", err)
	}
	got := string(bs)
	if got != msg {
		t.Errorf("response body want: %q, got %q", msg, got)
	}
}

func checkGet(t *testing.T, c *http.Client, url, msg string) {
	resp, err := c.Get(url)
	if err != nil {
		t.Fatalf("get %s: %s", url, msg)
	}
	defer resp.Body.Close()

	checkBody(t, resp, msg)
}

func TestServer(t *testing.T) {
	const msg = "hello"
	s, err := NewServer(
		[]string{"test.shanhu.io"},
		aries.StringFunc(msg),
	)
	if err != nil {
		t.Fatalf("create server: %s", err)
	}

	c := s.Client()
	resp, err := c.Get("https://test.shanhu.io")
	if err != nil {
		t.Fatalf("get: %s", err)
	}
	defer resp.Body.Close()

	checkBody(t, resp, msg)
}

func TestDualServer(t *testing.T) {
	const msg = "hello"

	s, err := NewDualServer(
		[]string{"test.shanhu.io"},
		aries.StringFunc(msg),
	)
	if err != nil {
		t.Fatalf("create server: %s", err)
	}

	c := s.Client()
	checkGet(t, c, "https://test.shanhu.io", msg)
	checkGet(t, c, "http://test.shanhu.io", msg)
}
