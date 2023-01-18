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

// Package srpc provides the Shanhu RPC caller for performing
// program-to-program, service-to-service or CLI-to-service interactions.
// It can also be used for Web-to-service AJAX-based RPC calls.
// The transport interface is intentionally limtied, and it is not RESTful.
// So this is not to be used to interact with a typical RESTful interface
// that is provided by third-party or other services.
package srpc

import (
	"net/url"

	"shanhu.io/pub/errcode"
)

// Caller is an RPC caller that can call
type Caller struct {
}

// NewCallerMust returns a caller where the address must be valid.
func NewCallerMust(addr string) *Caller {
	c, err := NewCaller(addr)
	if err != nil {
		panic(err)
	}
	return c
}

// NewCaller creates a new caller based on the given address.
func NewCaller(addr string) (*Caller, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, errcode.Annotate(err, "parse address")
	}
	return NewURLCaller(u), nil
}

// NewURLCaller returns a caller that calls to the specific URL.
func NewURLCaller(server *url.URL) *Caller {
	panic("TODO")
}
