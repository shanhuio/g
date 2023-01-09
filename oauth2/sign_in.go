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

type signInHandler struct {
	client   *Client
	redirect string
}

func newSignInHandler(client *Client, redirect string) *signInHandler {
	return &signInHandler{
		client:   client,
		redirect: redirect,
	}
}

func (h *signInHandler) Serve(c *aries.C) error {
	redirect := h.redirect
	if r := c.Req.URL.Query().Get("r"); r != "" {
		parsed, err := ParseRedirect(r)
		if err != nil {
			return err
		}
		redirect = parsed
	}
	state := &State{Dest: redirect}
	c.Redirect(h.client.SignInURL(state))
	return nil
}
