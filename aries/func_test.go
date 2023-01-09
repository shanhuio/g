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

	"shanhu.io/pub/httputil"
)

func TestFunc(t *testing.T) {
	const msg = "hello"
	f := StringFunc(msg)
	s := httptest.NewServer(Func(f))
	defer s.Close()

	got, err := httputil.GetString(s.Client(), s.URL)
	if err != nil {
		t.Error(err)
		return
	}
	if got != msg {
		t.Errorf("want %q in response, got %s", msg, got)
	}
}

func TestFuncHTTPS(t *testing.T) {
	const msg = "hello"
	f := StringFunc(msg)
	s := httptest.NewTLSServer(Func(f))
	defer s.Close()

	got, err := httputil.GetString(s.Client(), s.URL)
	if err != nil {
		t.Error(err)
		return
	}
	if got != msg {
		t.Errorf("want %q in response, got %s", msg, got)
	}
}
