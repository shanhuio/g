package creds

import (
	"context"
	"crypto/rsa"
	"net/http"
	"net/url"
	"sync"
	"time"

	"shanhu.io/pub/errcode"
	"shanhu.io/pub/httputil"
	"shanhu.io/pub/signer"
	"shanhu.io/pub/signin/signinapi"
	"shanhu.io/pub/timeutil"
)

type signInRequest struct {
	server *url.URL
	user   string
	key    *rsa.PrivateKey
}

func signIn(
	ctx context.Context, tr http.RoundTripper, req *signInRequest,
) (*signinapi.Creds, error) {
	signed, err := signer.RSASignTime(req.key)
	if err != nil {
		return nil, errcode.Annotate(err, "sign time")
	}

	const ttl = 30 * time.Minute

	c := &httputil.Client{Server: req.server, Transport: tr}
	sr := &signinapi.Request{
		User:        req.user,
		SignedTime:  signed,
		TTLDuration: timeutil.NewDuration(ttl),
	}
	sr.FillLegacyTTL()

	creds := new(signinapi.Creds)
	if err := c.CallContext(ctx, "/pubkey/signin", sr, creds); err != nil {
		return nil, errcode.Annotate(err, "sign in")
	}
	if creds.User != req.user {
		return nil, errcode.Internalf(
			"sign in as user %q, got %q", req.user, creds.User,
		)
	}

	creds.FixTime()

	return creds, nil
}

// TokenSource is a token source that fetches a signin token.
type TokenSource struct {
	req *signInRequest
}

// NewTokenSource returns a new token source that can get token from an account.
func NewTokenSource(
	server *url.URL, user string, key *rsa.PrivateKey,
) *TokenSource {
	req := &signInRequest{
		server: server,
		user:   user,
		key:    key,
	}
	return &TokenSource{req: req}
}

// Token gets the token.
func (s *TokenSource) Token(
	ctx context.Context, tr http.RoundTripper,
) (string, error) {
	creds, err := signIn(ctx, tr, s.req)
	if err != nil {
		return "", err
	}
	return creds.Token, nil
}

// CachingTokenSource is a token source that fetches a signin token,
// and caches the token until it expires.
type CachingTokenSource struct {
	req *signInRequest

	mu          sync.Mutex
	cached      *signinapi.Creds
	cacheExpire time.Time
}

// NewCachingTokenSource returns a new caching token source.
func NewCachingTokenSource(
	server *url.URL, user string, key *rsa.PrivateKey,
) *CachingTokenSource {
	req := &signInRequest{
		server: server,
		user:   user,
		key:    key,
	}
	return &CachingTokenSource{req: req}
}

func (s *CachingTokenSource) readCache() *signinapi.Creds {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cached == nil {
		return nil
	}

	now := time.Now()
	if now.Before(s.cacheExpire) {
		return s.cached
	}
	s.cached = nil
	return nil
}

func (s *CachingTokenSource) cache(creds *signinapi.Creds) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cached = creds
	s.cacheExpire = timeutil.Time(creds.ExpiresTime)
}

// Token gets a token. If there is one that is cached, and the cache has not
// expired, it returns the one that is cached.
func (s *CachingTokenSource) Token(
	ctx context.Context, tr http.RoundTripper,
) (string, error) {
	c := s.readCache()
	if c != nil {
		return c.Token, nil
	}

	c, err := signIn(ctx, tr, s.req)
	if err != nil {
		return "", err
	}
	s.cache(c)
	return c.Token, nil
}
