package creds

import (
	"context"
	"net/http"
	"net/url"

	"shanhu.io/g/httputil"
	"shanhu.io/g/osutil"
	"shanhu.io/g/rsautil"
	"shanhu.io/std/errcode"
)

// RobotEndpoint is an endpoint for robots.
type RobotEndpoint struct {
	Server    *url.URL
	User      string
	Key       []byte
	Transport http.RoundTripper
}

// LoadKeyFile loads the key from file f.
func (ep *RobotEndpoint) LoadKeyFile(f string) error {
	bs, err := osutil.ReadPrivateFile(f)
	if err != nil {
		return err
	}
	ep.Key = bs
	return nil
}

// Dial dials the server and returns a client with token.
func (ep *RobotEndpoint) Dial() (*httputil.Client, error) {
	k, err := rsautil.ParsePrivateKey(ep.Key)
	if err != nil {
		return nil, errcode.Annotate(err, "parse key")
	}

	tokenSrc := NewCachingTokenSource(ep.Server, ep.User, k)

	if _, err := tokenSrc.Token(context.TODO(), ep.Transport); err != nil {
		return nil, errcode.Annotate(err, "get token")
	}

	client := &httputil.Client{
		Server:      ep.Server,
		TokenSource: tokenSrc,
		Transport:   ep.Transport,
	}
	return client, nil
}
