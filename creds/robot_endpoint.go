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

package creds

import (
	"net/http"

	"shanhu.io/pub/errcode"
	"shanhu.io/pub/httputil"
	"shanhu.io/pub/osutil"
	"shanhu.io/pub/rsautil"
)

// RobotEndpoint is an endpoint for robots.
type RobotEndpoint struct {
	Server    string
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

	req := &Request{
		Server:    ep.Server,
		User:      ep.User,
		Key:       k,
		Transport: ep.Transport,
	}
	creds, err := NewCredsFromRequest(req)
	if err != nil {
		return nil, errcode.Annotate(err, "get creds")
	}
	client, err := httputil.NewTokenClient(ep.Server, creds.Token)
	if err != nil {
		return nil, err
	}
	if ep.Transport != nil {
		client.Transport = ep.Transport
	}
	return client, nil
}
