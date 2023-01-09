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
	"time"

	"shanhu.io/pub/aries"
	"shanhu.io/pub/errcode"
	"shanhu.io/pub/keyreg"
)

// JSONConfig is a JSON marshallable config that is commonly used for
// setting up a server.
type JSONConfig struct {
	GitHub       *App
	Google       *App
	DigitalOcean *App
	StateKey     string
	SessionKey   string
	SignInBypass string
	PublicKeys   map[string]string
}

// Config converts a JSON marshallable config to Config.
func (c *JSONConfig) Config() *Config {
	return &Config{
		GitHub:       c.GitHub,
		Google:       c.Google,
		DigitalOcean: c.DigitalOcean,
		StateKey:     []byte(c.StateKey),
		SessionKey:   []byte(c.SessionKey),
		Bypass:       c.SignInBypass,
		KeyRegistry:  keyreg.NewFileKeyRegistry(c.PublicKeys),
	}
}

// SimpleGitHubConfig converts a JSON marshallable config to Config that uses
// Github as the direct user ID mapping. Users that has a public key assigned
// in c.PublicKeys are defined as admin.
func (c *JSONConfig) SimpleGitHubConfig() *Config {
	ret := c.Config()
	ret.SignInCheck = MapGitHub
	ret.Check = func(name string) (interface{}, int, error) {
		if _, isAdmin := c.PublicKeys[name]; isAdmin {
			return nil, 1, nil
		}
		return nil, 0, nil
	}
	return ret
}

// Config is a module configuration for a GitHub Oauth handling module.
type Config struct {
	GitHub       *App
	Google       *App
	DigitalOcean *App

	StateKey        []byte
	SessionKey      []byte
	SessionLifeTime time.Duration
	SessionRefresh  time.Duration

	Bypass         string
	Redirect       string
	SignInRedirect string

	KeyRegistry keyreg.KeyRegistry

	// SignInCheck exchanges OAuth2 ID's for user ID.
	SignInCheck func(c *aries.C, u *UserMeta, purpose string) (string, error)

	// Check checks the user id and returns the user account structure.
	Check func(user string) (interface{}, int, error)

	PreSignOut func(c *aries.C) error
}

// MapGitHub is a login check function that only allows
// github login. It maps the user ID directly from GitHub users.
func MapGitHub(c *aries.C, u *UserMeta, _ string) (string, error) {
	if u.Method != MethodGitHub {
		return "", errcode.InvalidArgf(
			"login with %q not supported", u.Method,
		)
	}
	return u.ID, nil
}
