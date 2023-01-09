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

package signin

import (
	"time"

	"shanhu.io/pub/signin/signinapi"
	"shanhu.io/pub/timeutil"
)

// Token is a token with an expire time.
type Token struct {
	Token  string
	Expire time.Time
}

// Tokener issues auth tokens for users.
type Tokener interface {
	Token(user string, ttl time.Duration) *Token
}

// TokenCreds gets the credential from a token.
func TokenCreds(user string, tok *Token) *signinapi.Creds {
	return &signinapi.Creds{
		User:        user,
		Token:       tok.Token,
		ExpiresTime: timeutil.NewTimestamp(tok.Expire),
	}
}
