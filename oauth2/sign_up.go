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

package oauth2

import (
	"shanhu.io/pub/aries"
)

// SignUp is an HTTP module that handles user signups.
type SignUp struct {
	redirect string
	purpose  string
	signIn   bool
	module   *Module
	router   *aries.Router
}

// SignUpConfig is the config for creating a signup module.
type SignUpConfig struct {
	Redirect string

	// Whether keep user signed in after signing up.
	SignIn bool
}

// NewSignUp creates a new sign up module.
func NewSignUp(m *Module, c *SignUpConfig) *SignUp {
	s := &SignUp{
		purpose:  "signup",
		redirect: c.Redirect,
		module:   m,
		signIn:   c.SignIn,
	}

	s.router = s.makeRouter()
	return s
}

// Serve serves the incoming HTTP request.
func (s *SignUp) Serve(c *aries.C) error {
	return s.router.Serve(c)
}

func (s *SignUp) makeRouter() *aries.Router {
	r := aries.NewRouter()
	methods := s.module.Methods()
	for _, m := range methods {
		r.Get(m, s.handler(m))
	}
	return r
}

// Purpose returns the purpose string.
func (s *SignUp) Purpose() string {
	return s.purpose
}

func (s *SignUp) handler(m string) aries.Func {
	return func(c *aries.C) error {
		state := &State{
			Dest:     s.redirect,
			NoCookie: !s.signIn,
			Purpose:  s.purpose,
		}
		s.module.SignIn(c, m, state)
		return nil
	}
}
