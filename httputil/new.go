package httputil

import (
	"net/url"
)

// NewClientMust creates a client and panics on error.
func NewClientMust(s string) *Client { return NewTokenClientMust(s, "") }

// NewTokenClientMust creates a client with auth token and panics on error.
func NewTokenClientMust(s, tok string) *Client {
	c, err := NewTokenClient(s, tok)
	if err != nil {
		panic(err)
	}
	return c
}

// NewClient creates a new client.
func NewClient(s string) (*Client, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	return &Client{Server: u}, nil
}

// NewTokenClient creates a new client with a Bearer token.
func NewTokenClient(s, tok string) (*Client, error) {
	c, err := NewClient(s)
	if err != nil {
		return nil, err
	}
	c.TokenSource = NewStaticToken(tok)
	return c, nil
}

// NewUnixClient creates a new client that always goes to a particular
// unix domain socket.
func NewUnixClient(sockAddr string) *Client {
	return &Client{
		Server:    &url.URL{Scheme: "http", Host: "unix.sock"},
		Transport: unixSockTransport(sockAddr),
	}
}
