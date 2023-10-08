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
	"net/http"
	"net/url"
	"os"
	"os/user"

	"shanhu.io/g/errcode"
)

// Endpoint contains the login stub configuration.
type Endpoint struct {
	// Server is the server's prefix URL.
	Server *url.URL

	// User is an optional user name. If blank will use OS user name, or the
	// value of SHANHU_USER environment variable if exists.
	User string

	// Optional private key content. If nil, will use fall to use
	// PemFile. When presented, PemFile is ignored.
	Key []byte

	// Optional private key. If blank, will use the default key.
	PemFile string

	// Optional transport for creating the client.
	Transport http.RoundTripper

	Homeless bool // If true, will not look into the home folder for caches.
	NoTTY    bool // If true, will not fail if the key is encrypted.
}

// CurrentUser returns the new name of current user.
func CurrentUser() (string, error) {
	v, ok := os.LookupEnv("SHANHU_USER")
	if ok {
		return v, nil
	}

	u, err := user.Current()
	if err != nil {
		return "", err
	}
	return u.Username, nil
}

// NewEndpoint creates a new default endpoint for the target server.
func NewEndpoint(server string) (*Endpoint, error) {
	user, err := CurrentUser()
	if err != nil {
		return nil, err
	}
	u, err := url.Parse(server)
	if err != nil {
		return nil, errcode.Annotate(err, "parse URL")
	}
	return &Endpoint{User: user, Server: u}, nil
}

// NewRobot creates a new robot endpoint.
func NewRobot(
	server *url.URL, user, keyFile string, tr http.RoundTripper,
) *Endpoint {
	return &Endpoint{
		Server:    server,
		User:      user,
		PemFile:   keyFile,
		Homeless:  true,
		NoTTY:     true,
		Transport: tr,
	}
}
