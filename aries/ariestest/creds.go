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

package ariestest

import (
	"os"

	"shanhu.io/pub/creds"
	"shanhu.io/pub/errcode"
	"shanhu.io/pub/httputil"
)

// Login log into a server and fetch the token for the given user.
func Login(c *httputil.Client, user, key string) error {
	keyBytes, err := os.ReadFile(key)
	if err != nil {
		return errcode.Annotate(err, "read key")
	}
	endPoint := &creds.Endpoint{
		User:      user,
		Server:    c.Server.String(),
		Key:       keyBytes,
		Transport: c.Transport,
		Homeless:  true,
		NoTTY:     true,
	}

	login, err := creds.NewLogin(endPoint)
	if err != nil {
		return errcode.Annotate(err, "make login")
	}
	token, err := login.Token()
	if err != nil {
		return errcode.Annotate(err, "get token")
	}

	c.TokenSource = httputil.NewStaticToken(token)
	return nil
}
