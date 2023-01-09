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

package identity

import (
	"shanhu.io/pub/aries"
)

type service struct {
	card Card
}

func newService(card Card) *service {
	return &service{card: card}
}

// GetIDRequest is the request for getting an identity.
type GetIDRequest struct{}

func (s *service) apiGet(c *aries.C, req *GetIDRequest) (*Identity, error) {
	return s.card.Identity(c.Context)
}

// NewService creates a new identity service stub
func NewService(card Card) aries.Service {
	s := newService(card)
	r := aries.NewRouter()
	r.Call("get", s.apiGet)
	return r
}
