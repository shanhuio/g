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
	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
	"shanhu.io/g/aries"
	"shanhu.io/g/signer"
)

type digitalOcean struct{ c *Client }

var digitalOceanEndpoint = oauth2.Endpoint{
	AuthURL:  "https://cloud.digitalocean.com/v1/oauth/authorize",
	TokenURL: "https://cloud.digitalocean.com/v1/oauth/token",
}

func newDigitalOcean(
	app *App, s *signer.Sessions,
) *digitalOcean {
	c := NewClient(
		&oauth2.Config{
			ClientID:     app.ID,
			ClientSecret: app.Secret,
			Endpoint:     digitalOceanEndpoint,
		}, s, MethodDigitalOcean,
	)
	return &digitalOcean{c: c}
}

func (d *digitalOcean) client() *Client { return d.c }

func (d *digitalOcean) callback(c *aries.C) (*UserMeta, *State, error) {
	tok, state, err := d.c.TokenState(c)
	if err != nil {
		return nil, nil, err
	}

	oc := oauth2.NewClient(c.Context, oauth2.StaticTokenSource(tok))
	client := godo.NewClient(oc)

	account, _, err := client.Account.Get(c.Context)
	if err != nil {
		return nil, nil, err
	}

	return &UserMeta{
		Method: MethodDigitalOcean,
		ID:     account.UUID,
	}, state, nil
}
