package httputil

import (
	"context"
	"net/http"
)

// SetAuthToken sets authorization header token.
func SetAuthToken(h http.Header, tok string) {
	if tok == "" {
		return
	}
	h.Set("Authorization", "Bearer "+tok)
}

// TokenSource is an interface that can provides a bearer token for
// authentication.
type TokenSource interface {
	Token(ctx context.Context, tr http.RoundTripper) (string, error)
}

// StaticToken is a token source that provides a fixed,
// static token.
type StaticToken struct {
	T string
}

// NewStaticToken creates a new static token provider.
func NewStaticToken(tok string) *StaticToken {
	return &StaticToken{T: tok}
}

// Token always returns the fixed, static token T.
func (s *StaticToken) Token(
	ctx context.Context, tr http.RoundTripper,
) (string, error) {
	return s.T, nil
}
