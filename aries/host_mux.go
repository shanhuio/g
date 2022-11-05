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

package aries

// HostMux routes request to different services based on the incoming host.
type HostMux struct {
	m map[string]Service
}

// NewHostMux creates a new host mux.
func NewHostMux() *HostMux {
	return &HostMux{m: make(map[string]Service)}
}

// Set binds a host domain to a particular service.
func (m *HostMux) Set(host string, s Service) {
	m.m[host] = s
}

// Serve serves an incoming request.
func (m *HostMux) Serve(c *C) error {
	host := c.Req.Host
	s, found := m.m[host]
	if !found {
		return Miss
	}
	return s.Serve(c)
}
