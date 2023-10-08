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

package ariestest

import (
	"sort"

	"shanhu.io/g/aries"
	"shanhu.io/g/https/httpstest"
)

// HTTPSServers creates an HTTPS server that serves a set of
// domains.
func HTTPSServers(sites map[string]aries.Service) (*httpstest.Server, error) {
	m := aries.NewHostMux()
	var domains []string
	for domain, s := range sites {
		m.Set(domain, s)
		domains = append(domains, domain)
	}
	sort.Strings(domains)

	return httpstest.NewServer(domains, aries.Serve(m))
}

// HTTPSServer creates an HTTPS server that serves for a domain.
func HTTPSServer(domain string, s aries.Service) (*httpstest.Server, error) {
	return HTTPSServers(map[string]aries.Service{
		domain: s,
	})
}
