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

package sniproxy

import (
	"context"
)

// Router provides a host to connect with a token.
type Router interface {
	Route(ctx context.Context) (host string, token string, err error)
}

// StaticRouter routes to the given host with the given token.
type StaticRouter struct {
	Host  string
	Token string
}

// Route returns the given static host and token.
func (r *StaticRouter) Route(ctx context.Context) (string, string, error) {
	return r.Host, r.Token, nil
}
