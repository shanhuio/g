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
	"encoding/json"
	"fmt"

	"golang.org/x/oauth2"
	"shanhu.io/g/aries"
	"shanhu.io/g/signer"
	"shanhu.io/g/strutil"
)

type gitHub struct {
	c          *Client
	queryEmail bool
}

const gitHubEmailScope = "user:email"

// GitHubEndpoint is the oauth2 endpoint for GitHub.
var GitHubEndpoint = oauth2.Endpoint{
	AuthURL:  "https://github.com/login/oauth/authorize",
	TokenURL: "https://github.com/login/oauth/access_token",
}

func newGitHub(app *App, s *signer.Sessions) *gitHub {
	scopeSet := make(map[string]bool)
	if app.WithEmail {
		scopeSet[gitHubEmailScope] = true
	}
	for _, scope := range app.Scopes {
		scopeSet[scope] = true
	}
	scopes := strutil.SortedList(scopeSet)
	if scopes == nil {
		scopes = []string{}
	}
	queryEmail := scopeSet[gitHubEmailScope]

	c := NewClient(
		&oauth2.Config{
			ClientID:     app.ID,
			ClientSecret: app.Secret,
			Scopes:       scopes, // only need public information
			Endpoint:     GitHubEndpoint,
			RedirectURL:  app.RedirectURL,
		}, s, MethodGitHub,
	)
	return &gitHub{c: c, queryEmail: queryEmail}
}

func (g *gitHub) client() *Client { return g.c }

func (g *gitHub) callback(c *aries.C) (*UserMeta, *State, error) {
	tok, state, err := g.c.TokenState(c)
	if err != nil {
		return nil, nil, err
	}

	bs, err := getWithToken(c.Context, "https://api.github.com/user", tok)
	if err != nil {
		return nil, nil, err
	}

	var user struct {
		Login string `json:"login"`
		ID    int64  `json:"id"`
	}
	if err := json.Unmarshal(bs, &user); err != nil {
		return nil, nil, err
	}
	if user.ID == 0 {
		return nil, nil, fmt.Errorf("empty login")
	}

	var email string
	if g.queryEmail {
		const url = "https://api.github.com/user/emails"
		bs, err := getWithToken(c.Context, url, tok)
		if err != nil {
			return nil, nil, err
		}

		type userEmail struct {
			Email    string `json:"email"`
			Verified bool   `json:"verified"`
			Primary  bool   `json:"primary"`
		}

		var emails []*userEmail
		if err := json.Unmarshal(bs, &emails); err != nil {
			return nil, nil, err
		}
		for _, m := range emails {
			if m.Primary && m.Verified {
				email = m.Email
			}
		}
	}
	meta := &UserMeta{
		Method: MethodGitHub,
		ID:     fmt.Sprint(user.ID),
		Name:   user.Login,
		Email:  email,
	}
	return meta, state, nil
}
