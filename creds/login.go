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
	"time"

	"shanhu.io/pub/errcode"
	"shanhu.io/pub/httputil"
	"shanhu.io/pub/timeutil"
)

// LoginWithKey uses the given PEM file to login a server, and returns the creds
// if succeess.
func LoginWithKey(p *Endpoint) (*Creds, error) {
	k, err := readEndpointKey(p)
	if err != nil {
		return nil, err
	}
	req := &Request{
		Server:    p.Server,
		User:      p.User,
		Key:       k,
		Transport: p.Transport,
	}
	return NewCredsFromRequest(req)
}

// Login is a helper stub to perform login actions.
type Login struct {
	endPoint   *Endpoint
	credsStore credsStore
	creds      *Creds // cached creds
}

// NewServerLogin returns a new server login with default user and pem file.
func NewServerLogin(s string) (*Login, error) {
	p, err := NewEndpoint(s)
	if err != nil {
		return nil, err
	}
	return NewLogin(p)
}

// NewLogin creates a new login stub with the given config.
func NewLogin(p *Endpoint) (*Login, error) {
	if p.User == "" {
		return nil, errcode.InvalidArgf("user is empty")
	}

	cp := *p
	if cp.PemFile == "" && !cp.Homeless {
		pem, err := HomeFile("key.pem")
		if err != nil {
			return nil, errcode.Internalf("fail to get home: %v", err)
		}
		cp.PemFile = pem
	}

	lg := &Login{endPoint: &cp}
	if !p.Homeless {
		lg.credsStore = newHomeCredsStore(p.Server.String())
	}
	return lg, nil
}

// NewRobotLogin is a shorthand for NewLogin(NewRobot())
func NewRobotLogin(
	server *url.URL, user, keyFile string, tr http.RoundTripper,
) (*Login, error) {
	return NewLogin(NewRobot(server, user, keyFile, tr))
}

func (lg *Login) check(cs *Creds) (bool, error) {
	if cs.User != lg.endPoint.User {
		return false, nil
	}
	if cs.Server != lg.endPoint.Server.String() {
		return false, nil
	}

	expires := timeutil.Time(cs.Creds.ExpiresTime)
	now := time.Now()
	if !now.Before(expires) {
		return false, nil
	}

	return true, nil
}

// Token returns the login token for the login. If a valid token is already
// cached, it returns the cached one.
func (lg *Login) Token() (string, error) {
	if lg.endPoint.Homeless {
		// Nothing cached anywhere, just return a new one.
		return lg.GetToken()
	}

	cs := lg.creds
	if cs == nil && lg.credsStore != nil {
		newCreds, err := lg.credsStore.read()
		if err != nil {
			if errcode.IsNotFound(err) {
				return lg.GetToken()
			}
			return "", err
		}
		if newCreds == nil {
			panic("should have creds loaded from the file system")
		}
		newCreds.FixTime()
		cs = newCreds
	}

	// now we loaded a cached creds
	ok, err := lg.check(cs)
	if err != nil {
		return "", err
	}
	if !ok {
		return lg.GetToken()
	}

	return cs.Token, nil
}

// Do performs the login and returns the credentials.
// It does not read or write the credential cache file.
func (lg *Login) Do() (*Creds, error) {
	return LoginWithKey(lg.endPoint)
}

// GetToken returns the login token for the login. It ignores and overwrites
// any existing login token that uses the same login creds file.
func (lg *Login) GetToken() (string, error) {
	cs, err := lg.Do()
	if err != nil {
		return "", err
	}

	// cache it
	lg.creds = cs

	// If not homeless, also cache it in home directory.
	if lg.credsStore != nil {
		if err := lg.credsStore.write(cs); err != nil {
			return "", err
		}
	}
	return cs.Creds.Token, nil
}

// Dial creates an token client.
func (lg *Login) Dial() (*httputil.Client, error) {
	tok, err := lg.Token()
	if err != nil {
		return nil, err
	}

	c := &httputil.Client{
		Server:      lg.endPoint.Server,
		TokenSource: httputil.NewStaticToken(tok),
	}
	c.Transport = lg.endPoint.Transport
	return c, nil
}

type loginTokenSource struct {
	login *Login
}

func (s *loginTokenSource) Token(
	_ context.Context, tr http.RoundTripper,
) (string, error) {
	return s.login.Token()
}

// TokenSource converts the login to a TokenSource.
func (lg *Login) TokenSource() httputil.TokenSource {
	return &loginTokenSource{login: lg}
}

// LoginServer uses the default setting to login into a server.
func LoginServer(server string) (string, error) {
	login, err := NewServerLogin(server)
	if err != nil {
		return "", err
	}
	return login.Token()
}

// Dial logins the server and returns the httputil client.
func Dial(server string) (*httputil.Client, error) {
	tok, err := LoginServer(server)
	if err != nil {
		return nil, err
	}
	return httputil.NewTokenClient(server, tok)
}

// DialAsUser logins the server as user and returns the httputil client.
func DialAsUser(user, server string) (*httputil.Client, error) {
	if user == "" {
		return Dial(server)
	}
	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, errcode.Annotate(err, "parse server")
	}

	ep := &Endpoint{
		User:   user,
		Server: serverURL,
	}
	login, err := NewLogin(ep)
	if err != nil {
		return nil, err
	}
	tok, err := login.Token()
	if err != nil {
		return nil, err
	}
	return httputil.NewTokenClient(server, tok)
}

// DialEndpoint creates a token client with the given endpoint.
func DialEndpoint(p *Endpoint) (*httputil.Client, error) {
	login, err := NewLogin(p)
	if err != nil {
		return nil, err
	}
	return login.Dial()
}
