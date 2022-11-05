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

// Package redhttp provides a service that routes
// all incoming http requests to https.
package redhttp

import (
	"flag"
	"log"

	"shanhu.io/pub/aries"
	"shanhu.io/pub/errcode"
)

// Redirect redirects the incoming request to https.
func Redirect(c *aries.C) error {
	host := c.Req.Host
	if host == "" {
		return errcode.InvalidArgf("host not found")
	}
	u := *c.Req.URL
	u.Host = host
	u.Scheme = "https"
	c.Redirect(u.String())
	return nil
}

// Main is the main entrance for the service.
func Main() {
	addr := aries.DeclareAddrFlag("localhost:8000")
	flag.Parse()
	if err := aries.ListenAndServe(*addr, aries.Func(Redirect)); err != nil {
		log.Fatal(err)
	}
}
