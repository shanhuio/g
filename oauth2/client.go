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
	"context"
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
	"shanhu.io/pub/aries"
	"shanhu.io/pub/errcode"
	"shanhu.io/pub/httputil"
	"shanhu.io/pub/signer"
)

func stateCode(req *http.Request) (state, code string) {
	values := req.URL.Query()
	state = values.Get("state")
	if state != "" {
		code = values.Get("code")
	}
	return state, code
}

// State contains a JSON marshalable state for OAuth2 sign in.
type State struct {
	// URL to redirect to after signing in.
	Dest string

	// Sign in purpose.
	Purpose string `json:",omitempty"`

	// Whether set cookie after signing in.
	NoCookie bool `json:",omitempty"`
}

// Client is an oauth client for oauth2 exchanges.
type Client struct {
	config *oauth2.Config
	states *signer.Sessions
	method string
}

// NewClient creates a new oauth client for oauth2 exchnages.
func NewClient(
	c *oauth2.Config, states *signer.Sessions, m string,
) *Client {
	return &Client{
		config: c,
		states: states,
		method: m,
	}
}

// Method returns the method class of this oauth2 client.
func (c *Client) Method() string { return c.method }

// SignInURL returns the online signin URL for redirection.
func (c *Client) SignInURL(s *State) string {
	state, _, err := c.states.NewJSON(s)
	if err != nil {
		panic(err)
	}
	return c.config.AuthCodeURL(state)
}

// OfflineSignInURL returns the offline signin URL for redirection.
func (c *Client) OfflineSignInURL(s *State) string {
	state, _, err := c.states.NewJSON(s)
	if err != nil {
		panic(err)
	}
	return c.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// TokenState extracts the oauth2 access token and state from the request.
func (c *Client) TokenState(ctx *aries.C) (*oauth2.Token, *State, error) {
	stateStr, code := stateCode(ctx.Req)
	if stateStr == "" {
		return nil, nil, fmt.Errorf("invalid oauth redirect")
	}

	state := new(State)
	if !c.states.CheckJSON(stateStr, state) {
		return nil, nil, fmt.Errorf("state invalid")
	}

	tok, err := c.config.Exchange(ctx.Context, code)
	if err != nil {
		return nil, nil, fmt.Errorf("exchange failed: %v", err)
	}
	if !tok.Valid() {
		return nil, nil, fmt.Errorf("token is invalid")
	}
	return tok, state, nil
}

func getWithToken(
	ctx context.Context, u string, tok *oauth2.Token,
) ([]byte, error) {
	urlParsed, err := url.Parse(u)
	if err != nil {
		return nil, errcode.Annotate(err, "parse URL")
	}

	serverURL := &url.URL{
		Scheme: urlParsed.Scheme,
		Opaque: urlParsed.Opaque,
		Host:   urlParsed.Host,
	}

	client := &httputil.Client{
		Server:      serverURL,
		TokenSource: httputil.NewStaticToken(tok.AccessToken),
	}
	return client.GetBytes(urlParsed.Path)
}
