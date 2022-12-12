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
	"context"
	"crypto/rsa"
	"net/http"
	"net/url"
	"time"

	"shanhu.io/pub/signin/signinapi"
)

// Creds is the credential that is cached after logging in. This can also be
// saved in JSON format in user's home directory.
type Creds struct {
	Server          string
	signinapi.Creds // User name is saved in this.
}

// Request contains the configuration to create a credential.
type Request struct {
	Server *url.URL
	User   string
	Key    *rsa.PrivateKey
	TTL    time.Duration

	// Transport is the http transport for the token exchange.
	Transport http.RoundTripper
}

// NewCredsFromRequest creates a new user credential by dialing the server
// using the given RSA private key.
func NewCredsFromRequest(req *Request) (*Creds, error) {
	sr := &signInRequest{
		server: req.Server,
		user:   req.User,
		key:    req.Key,
	}

	creds, err := signIn(context.TODO(), req.Transport, sr)
	if err != nil {
		return nil, err
	}

	return &Creds{
		Creds:  *creds,
		Server: req.Server.String(),
	}, nil
}

// NewCreds creates a new user credential by dialing the server using
// the given RSA private key.
func NewCreds(server *url.URL, user string, k *rsa.PrivateKey) (*Creds, error) {
	req := &Request{
		Server: server,
		User:   user,
		Key:    k,
	}
	return NewCredsFromRequest(req)
}
