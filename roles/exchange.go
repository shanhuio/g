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
		req.User, req.SelfToken, t,
	); err != nil {
		return nil, altAuthErr(err, "verify self token")
	}

	const ttl = 30 * time.Minute
	token := x.tokener.Token(req.User, ttl)
	return signin.TokenCreds(req.User, token), nil
}
