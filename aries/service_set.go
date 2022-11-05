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

// ServiceSet is a set of muxes that
type ServiceSet struct {
	Auth Auth

	Resource Service
	Guest    Service
	User     Service
	Admin    Service
	IsAdmin  func(c *C) bool

	InternalSignIn Func
}

func serveService(m Service, c *C) error {
	if m == nil {
		return Miss
	}
	return m.Serve(c)
}

func (s *ServiceSet) isAdmin(c *C) bool {
	if s.IsAdmin == nil {
		return c.User != "" && c.UserLevel > 0
	}
	return s.IsAdmin(c)
}

// serveAuth performs the auth check.
func (s *ServiceSet) serveAuth(c *C) (bool, error) {
	if err := s.Auth.Serve(c); err != Miss {
		return true, err
	}
	return false, s.Auth.Setup(c)
}

// Serve serves the incoming request with the mux set.
func (s *ServiceSet) Serve(c *C) error {
	if served, err := s.serveAuth(c); err != nil {
		return err
	} else if served {
		return nil
	}

	if err := serveService(s.Resource, c); err != Miss {
		return err
	}
	if err := serveService(s.Guest, c); err != Miss {
		return err
	}
	if c.User != "" {
		if err := serveService(s.User, c); err != Miss {
			return err
		}
	}
	if s.isAdmin(c) {
		if err := serveService(s.Admin, c); err != Miss {
			return err
		}
	}

	return Miss
}

// ServeInternal serves the incoming request with the mux set, but only serves
// resource for normal users, and allows only admins (users with positive
// level) to visit the guest mux.
func (s *ServiceSet) ServeInternal(c *C) error {
	if err := serveService(s.Auth, c); err != Miss {
		return err
	}
	if s.Auth != nil {
		if err := s.Auth.Setup(c); err != nil {
			return err
		}
	}

	if err := serveService(s.Resource, c); err != Miss {
		return err
	}

	if !s.isAdmin(c) {
		if c.Path == "/" {
			if s.InternalSignIn != nil {
				return s.InternalSignIn(c)
			}
			return NeedSignIn
		}
		c.Redirect("/")
		return nil
	}

	if err := serveService(s.Guest, c); err != Miss {
		return err
	}
	if err := serveService(s.User, c); err != Miss {
		return err
	}
	if err := serveService(s.Admin, c); err != Miss {
		return err
	}

	return Miss
}
