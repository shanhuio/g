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

package redirect

import (
	"flag"
	"log"

	"shanhu.io/g/aries"
)

type config struct {
	RedirectToDomain string
}

type server struct {
	c *config
}

func (s *server) redirect(c *aries.C) error {
	u := *c.Req.URL // make a shallow copy
	u.Scheme = "https"
	u.Host = s.c.RedirectToDomain
	c.Redirect(u.String())
	return nil
}

func newServer(c *config) (aries.Func, error) {
	s := &server{c: c}
	return s.redirect, nil
}

// Main is the main entrance for the redirect service.
func Main() {
	addr := aries.DeclareAddrFlag("localhost:8000")
	to := flag.String("to", "", "redirect to this address")
	flag.Parse()

	config := &config{RedirectToDomain: *to}

	f, err := newServer(config)
	if err != nil {
		log.Fatal(err)
	}
	if err := aries.ListenAndServe(*addr, f); err != nil {
		log.Fatal(err)
	}
}
