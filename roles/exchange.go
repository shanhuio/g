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

package roles

import (
	"time"

	"shanhu.io/pub/aries"
	"shanhu.io/pub/signin"
	"shanhu.io/pub/signin/signinapi"
)

// SignInRequest signs in with a self-signed token.
type SignInRequest struct {
	User      string
	SelfToken string
}

// Exchange is a token exchange that exchanges self token
// for session token and optionally
type Exchange struct {
	roles   *Roles
	tokener signin.Tokener
}

// NewExchange creates a new exchange that can exchange self token of r into
// access tokens issued by tokener.
func NewExchange(r *Roles, tokener signin.Tokener) *Exchange {
	return &Exchange{
		roles:   r,
		tokener: tokener,
	}
}

// Exchange exchanges self token for session token.
func (x *Exchange) Exchange(c *aries.C, req *SignInRequest) (
	*signinapi.Creds, error,
) {
	t := time.Now()
	if _, err := x.roles.VerifySelfToken(
		c.Context, req.User, req.SelfToken, t,
	); err != nil {
		return nil, altAuthErr(err, "verify self token")
	}

	const ttl = 30 * time.Minute
	token := x.tokener.Token(req.User, ttl)
	return signin.TokenCreds(req.User, token), nil
}
