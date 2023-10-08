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

package authgate

import (
	"time"

	"shanhu.io/g/aries"
	"shanhu.io/g/errcode"
	"shanhu.io/g/identity"
	"shanhu.io/g/jwt"
	"shanhu.io/g/signin"
	"shanhu.io/g/signin/signinapi"
	"shanhu.io/g/timeutil"
)

// ExchangeConfig is the config for creating an session exchanger
// that exchanges access tokens for session tokens.
type ExchangeConfig struct {
	Audience string
	Issuer   string
	Card     identity.Card
	Now      func() time.Time
}

// Exchange exchanges an access tokens for a session token. An access token is
// a JWT that is signed by a realm CA, as a proof that the client has been
// authorized to access some resource on behalf of the user for a period of
// time. The session token is a token that is issued by a local tokener, which
// can be used to access the API. Checking a session token is often a much
// light-weight local operation, which does not require querying the central
// realm.
type Exchange struct {
	audience string
	issuer   string
	card     identity.Card
	verifier jwt.Verifier
	tokener  signin.Tokener
	now      func() time.Time
}

// NewExchange creates an exchange that exchnages access tokens
// for session tokens from tok.
func NewExchange(tok signin.Tokener, config *ExchangeConfig) *Exchange {
	return &Exchange{
		audience: config.Audience,
		issuer:   config.Issuer,
		card:     config.Card,
		verifier: identity.NewJWTVerifier(config.Card),
		tokener:  tok,
		now:      timeutil.NowFunc(config.Now),
	}
}

// Exchange is the API that exchanges access tokens for session tokens in the
// form of credentials.
func (x *Exchange) Exchange(c *aries.C, req *signinapi.Request) (
	*signinapi.Creds, error,
) {
	if req.AccessToken == "" {
		return nil, errcode.InvalidArgf("access token missing")
	}

	ctx := c.Context

	now := x.now()
	tok, err := jwt.DecodeAndVerify(ctx, req.AccessToken, x.verifier, now)
	if err != nil {
		return nil, errcode.Annotate(err, "invalid token")
	}

	wantClaims := &jwt.ClaimSet{
		Sub: req.User,
		Iss: x.issuer,
		Aud: x.audience,
	}
	if err := jwt.CheckClaimSet(tok.ClaimSet, wantClaims); err != nil {
		return nil, errcode.Annotate(err, "invalid claims")
	}

	ttl := timeutil.TimeDuration(req.TTLDuration)
	if ttl <= time.Duration(0) {
		return nil, errcode.Unauthorizedf("ttl too short")
	}

	token := x.tokener.Token(req.User, ttl)
	return signin.TokenCreds(req.User, token), nil
}
