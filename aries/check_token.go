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

package aries

import (
	"fmt"
	"strings"

	"shanhu.io/pub/signer"
)

// Bearer returns the authorization token.
func Bearer(c *C) string {
	auth := c.Req.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return ""
	}
	return strings.TrimPrefix(auth, "Bearer ")
}

// CheckToken checks if the bearer token is properly signed by the
// same API key.
func CheckToken(c *C, s *signer.TimeSigner) error {
	token := Bearer(c)
	if !s.Check(token) {
		return fmt.Errorf("invalid token")
	}

	return nil
}
