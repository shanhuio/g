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

import (
	"net/http"
)

// Service is a interface similar to Func
type Service interface {
	Serve(c *C) error
}

// Serve wraps a service into an HTTP handler.
func Serve(s Service) http.Handler {
	f, ok := s.(Func)
	if ok {
		// if it is a function already, we do not need to do the wrapping.
		return f
	}
	return Func(s.Serve)
}
