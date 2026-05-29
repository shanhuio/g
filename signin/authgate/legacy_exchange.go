package authgate

import (
	"time"

	"shanhu.io/g/aries"
	"shanhu.io/g/keyreg"
	"shanhu.io/g/signer"
	"shanhu.io/g/signin"
	"shanhu.io/g/signin/signinapi"
	"shanhu.io/std/errcode"
)

// LegacyExchange handles sign in using a public key registry. The request
// presents a signed time using the user's private key to authenticate.
type LegacyExchange struct {
	tokener     signin.Tokener
	keyRegistry keyreg.KeyRegistry
}

// NewLegacyExchange creates a legacy public key based credential exchange
// where the client presents a signed time with its private key.
func NewLegacyExchange(
	tok signin.Tokener, reg keyreg.KeyRegistry,
) *LegacyExchange {
	return &LegacyExchange{
		tokener:     tok,
		keyRegistry: reg,
	}
}

// Exchange handles the request to exchange a public-key signed timestamp to a
// token.
func (x *LegacyExchange) Exchange(c *aries.C, req *signinapi.Request) (
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
