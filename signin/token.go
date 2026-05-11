package signin

import (
	"time"

	"shanhu.io/g/signin/signinapi"
	"shanhu.io/g/timeutil"
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
