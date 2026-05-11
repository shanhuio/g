package creds

import (
	"context"
	"crypto/rsa"
	"net/http"
	"net/url"
	"time"

	"shanhu.io/g/signin/signinapi"
)

// Creds is the credential that is cached after logging in. This can also be
// saved in JSON format in user's home directory.
type Creds struct {
	Server          string
	signinapi.Creds // User name is saved in this.
}

// Request contains the configuration to create a credential.
type Request struct {
	Server *url.URL
	User   string
	Key    *rsa.PrivateKey
	TTL    time.Duration

	// Transport is the http transport for the token exchange.
	Transport http.RoundTripper
}

// NewCredsFromRequest creates a new user credential by dialing the server
// using the given RSA private key.
func NewCredsFromRequest(req *Request) (*Creds, error) {
	sr := &signInRequest{
		server: req.Server,
		user:   req.User,
		key:    req.Key,
	}

	creds, err := signIn(context.TODO(), req.Transport, sr)
	if err != nil {
		return nil, err
	}

	return &Creds{
		Creds:  *creds,
		Server: req.Server.String(),
	}, nil
}

// NewCreds creates a new user credential by dialing the server using
// the given RSA private key.
func NewCreds(server *url.URL, user string, k *rsa.PrivateKey) (*Creds, error) {
	req := &Request{
		Server: server,
		User:   user,
		Key:    k,
	}
	return NewCredsFromRequest(req)
}
