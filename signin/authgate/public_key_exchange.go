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

	"shanhu.io/pub/aries"
	"shanhu.io/pub/errcode"
	"shanhu.io/pub/keyreg"
	"shanhu.io/pub/signer"
	"shanhu.io/pub/signin"
	"shanhu.io/pub/signin/signinapi"
)

// PublicKeyExchange handles sign in using a public key registry. The request
// presents a signed time using the user's private key to authenticate.
type PublicKeyExchange struct {
	tokener     signin.Tokener
	keyRegistry keyreg.KeyRegistry
}

// NewPublicKeyExchange creates a legacy public key based credential exchange
// where the client presents a signed time with its private key.
func NewPublicKeyExchange(
	tok signin.Tokener, reg keyreg.KeyRegistry,
) *PublicKeyExchange {
	return &PublicKeyExchange{
		tokener:     tok,
		keyRegistry: reg,
	}
}

// Exchange handles the request to exchange a public-key signed timestamp to a
// token.
func (x *PublicKeyExchange) Exchange(c *aries.C, req *signinapi.Request) (
	*signinapi.Creds, error,
) {
	if req.SignedTime == nil {
		return nil, errcode.InvalidArgf("signature missing")
	}

	keys, err := x.keyRegistry.Keys(req.User)
	if err != nil {
		return nil, err
	}

	key := keyreg.FindKeyByHash(keys, req.SignedTime.KeyID)
	if key == nil {
		return nil, errcode.Unauthorizedf("signing key not authorized")
	}

	const window = time.Minute * 5
	if err := signer.CheckRSATimeSignature(
		req.SignedTime, key.Key(), window,
	); err != nil {
		return nil, errcode.Add(errcode.Unauthorized, err)
	}

	ttl := req.GetTTL()
	token := x.tokener.Token(req.User, ttl)
	return signin.TokenCreds(req.User, token), nil
}
