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

package signin

import (
	"time"

	"shanhu.io/pub/aries"
	"shanhu.io/pub/signer"
	"shanhu.io/pub/timeutil"
)

const cookieKey = "session"

func defaultCheck(user string) (interface{}, int, error) {
	lvl := 0
	if user != "" {
		lvl = 1
	}
	return nil, lvl, nil
}

// GateConfig contains configuration for initializing an
// identity gate.
type GateConfig struct {
	SessionKey      []byte
	SessionLifeTime time.Duration
	SessionRefresh  time.Duration

	Check func(user string) (interface{}, int, error)
}

// Gate is a token checking gate that checks the auth token
// and saves the checking result in the context.
type Gate struct {
	sessions       *signer.Sessions
	sessionRefresh time.Duration

	check func(user string) (interface{}, int, error)
}

// NewGate creates a new auth token checking gate.
func NewGate(config *GateConfig) *Gate {
	sessionLifeTime := config.SessionLifeTime
	if sessionLifeTime <= 0 {
		sessionLifeTime = timeutil.Week
	}
	sessionRefresh := config.SessionRefresh
	if sessionRefresh <= 0 || sessionRefresh > sessionLifeTime {
		sessionRefresh = sessionLifeTime / 5 * 4
	}
	sessions := signer.NewSessions(config.SessionKey, sessionLifeTime)

	check := config.Check
	if check == nil {
		check = defaultCheck
	}

	return &Gate{
		sessions:       sessions,
		sessionRefresh: sessionRefresh,
		check:          check,
	}
}

// Different token types.
const (
	TokenCookie = "cookie"
	TokenBearer = "bearer"
)

func authToken(c *aries.C) (string, string) {
	if bearer := aries.Bearer(c); bearer != "" {
		return bearer, TokenBearer
	}
	return c.ReadCookie(cookieKey), TokenCookie
}

// CheckToken checks a token of a specific type.
func (g *Gate) CheckToken(token, typ string) (*CredsInfo, error) {
	info := &CredsInfo{TokenType: typ}
	bs, left, ok := g.sessions.Check(token)
	if !ok {
		return info, nil
	}
	info.NeedRefresh = left < g.sessionRefresh

	user := string(bs)
	dat, lvl, err := g.check(user)
	if err != nil {
		return nil, err
	}
	info.User = user
	info.UserLevel = lvl
	info.Valid = lvl >= 0
	info.Data = dat

	return info, nil
}

// Check checks the auth token in the context, with the session
// signature and the check function. It returns if it is valid, if it needs
// refresh.
func (g *Gate) Check(c *aries.C) (*CredsInfo, error) {
	return g.CheckToken(authToken(c))
}

// Token returns an auth token that is valid for ttl. It returns the token
// and the expiry time.
func (g *Gate) Token(user string, ttl time.Duration) *Token {
	token, expire := g.sessions.New([]byte(user), ttl)
	return &Token{
		Token:  token,
		Expire: expire,
	}
}

// SetupCookie sets up the cookie for a particular user.
func (g *Gate) SetupCookie(c *aries.C, user string) {
	token := g.Token(user, 0)
	c.WriteCookie(cookieKey, token.Token, token.Expire)
}

// ClearGateCookie clears the gate's session cookie.
func ClearGateCookie(c *aries.C) {
	c.ClearCookie(cookieKey)
}

// CheckAndSetup checks the user credentials. If the credential is valid it
// also applies the credential to the context. If the credential is not
// valid, it clears the cookie. If the credential needs refreshing
// it refreshes the cookie.
func (g *Gate) CheckAndSetup(c *aries.C) (bool, error) {
	creds, err := g.Check(c)
	if err != nil {
		return false, err
	}

	ApplyCredsInfo(c, creds)

	if creds.TokenType == TokenCookie {
		if !creds.Valid {
			ClearGateCookie(c)
		} else if creds.NeedRefresh {
			g.SetupCookie(c, creds.User)
		}
	}

	return creds.Valid, nil
}

// Setup sets up the credentials for the request.
func (g *Gate) Setup(c *aries.C) error {
	_, err := g.CheckAndSetup(c)
	return err
}

// Serve serves nothing. It is defined just to satisfy aries.Auth interface.
func (g *Gate) Serve(c *aries.C) error {
	return aries.Miss // A simple gate does not serve anything.
}
