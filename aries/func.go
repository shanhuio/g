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
	"io"
	"net"
	"net/http"
)

// Func defines an HTTP handling function.
type Func func(c *C) error

// Serve implements the service interface.
func (f Func) Serve(c *C) error { return f(c) }

func (f Func) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := NewContext(w, req)
	c.ErrCode(f(c))
}

// ListenAndServe launches the handler as an HTTP service.
func (f Func) ListenAndServe(addr string) error {
	s := &http.Server{
		Addr:    addr,
		Handler: f,
	}
	return s.ListenAndServe()
}

// ServeAt launches the handler as an HTTP service at the given
// listener.
func (f Func) ServeAt(lis net.Listener) error {
	s := &http.Server{Handler: f}
	return s.Serve(lis)
}

// StringFunc creates a Func that always reply the given string.
func StringFunc(s string) Func {
	return func(c *C) error {
		io.WriteString(c.Resp, s)
		return nil
	}
}

// RedirectTo creates a Func that always redirects to u
func RedirectTo(u string) Func {
	return func(c *C) error {
		c.Redirect(u)
		return nil
	}
}
