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

package oauth2

import (
	"context"
	"encoding/json"
	"fmt"

	"golang.org/x/oauth2"
	"shanhu.io/pub/aries"
	"shanhu.io/pub/signer"
	"shanhu.io/pub/strutil"
)

// GoogleUserInfo stores a Google user's basic personal info.
type GoogleUserInfo struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// GetGoogleUserInfo queries Google OAuth endpoint for user info data.
func GetGoogleUserInfo(
	ctx context.Context, c *Client, tok *oauth2.Token,
) (*GoogleUserInfo, error) {
	const apiPath = "https://www.googleapis.com/oauth2/v3/userinfo"
	bs, err := getWithToken(ctx, apiPath, tok)
	if err != nil {
		return nil, err
	}

	user := new(GoogleUserInfo)
	if err := json.Unmarshal(bs, user); err != nil {
		return nil, err
	}

	return user, nil
}

const (
	googleEmailScope   = "https://www.googleapis.com/auth/userinfo.email"
	googleProfileScope = "https://www.googleapis.com/auth/userinfo.profile"
)

// GoogleEndpoint is google's oauth2 endpoint. This is copied from
// golang.org/x/oauth2/google, to avoid dragging in useless dependencies such
// as protobuf.
//
// Google engineers really should be more careful with their package dependency
// management.
var GoogleEndpoint = oauth2.Endpoint{
	AuthURL:   "https://accounts.google.com/o/oauth2/auth",
	TokenURL:  "https://oauth2.googleapis.com/token",
	AuthStyle: oauth2.AuthStyleInParams,
}

type google struct{ c *Client }

func newGoogle(app *App, s *signer.Sessions) *google {
	scopeSet := make(map[string]bool)
	// Google OAuth has to have at least one scope to get user ID.
	scopeSet[googleEmailScope] = true
	if app.WithProfile {
		scopeSet[googleProfileScope] = true
	}
	scopes := strutil.SortedList(scopeSet)
	if scopes == nil {
		scopes = []string{}
	}
	c := NewClient(
		&oauth2.Config{
			ClientID:     app.ID,
			ClientSecret: app.Secret,
			Scopes:       scopes,
			Endpoint:     GoogleEndpoint,
			RedirectURL:  app.RedirectURL,
		}, s, MethodGoogle,
	)
	return &google{c: c}
}

func (g *google) client() *Client { return g.c }

func (g *google) callback(c *aries.C) (*UserMeta, *State, error) {
	tok, state, err := g.c.TokenState(c)
	if err != nil {
		return nil, nil, err
	}

	user, err := GetGoogleUserInfo(c.Context, g.c, tok)
	if err != nil {
		return nil, nil, err
	}

	email := user.Email
	if email == "" {
		return nil, nil, fmt.Errorf("empty login")
	}
	name := user.Name
	if name == "" {
		name = "no-name"
	}
	return &UserMeta{
		Method: MethodGoogle,
		ID:     email,
		Name:   name,
		Email:  email,
	}, state, nil
}
