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

package creds

import (
	"context"
	"net/http"
	"net/url"

	"shanhu.io/g/errcode"
	"shanhu.io/g/httputil"
	"shanhu.io/g/osutil"
	"shanhu.io/g/rsautil"
)

// RobotEndpoint is an endpoint for robots.
type RobotEndpoint struct {
	Server    *url.URL
	User      string
	Key       []byte
	Transport http.RoundTripper
}

// LoadKeyFile loads the key from file f.
func (ep *RobotEndpoint) LoadKeyFile(f string) error {
	bs, err := osutil.ReadPrivateFile(f)
	if err != nil {
		return err
	}
	ep.Key = bs
	return nil
}

// Dial dials the server and returns a client with token.
func (ep *RobotEndpoint) Dial() (*httputil.Client, error) {
	k, err := rsautil.ParsePrivateKey(ep.Key)
	if err != nil {
		return nil, errcode.Annotate(err, "parse key")
	}

	tokenSrc := NewCachingTokenSource(ep.Server, ep.User, k)

	if _, err := tokenSrc.Token(context.TODO(), ep.Transport); err != nil {
		return nil, errcode.Annotate(err, "get token")
	}

	client := &httputil.Client{
		Server:      ep.Server,
		TokenSource: tokenSrc,
		Transport:   ep.Transport,
	}
	return client, nil
}
