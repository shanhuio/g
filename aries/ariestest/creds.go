package ariestest

import (
	"os"

	"shanhu.io/g/creds"
	"shanhu.io/g/httputil"
	"shanhu.io/std/errcode"
)

// Login log into a server and fetch the token for the given user.
func Login(c *httputil.Client, user, key string) error {
	keyBytes, err := os.ReadFile(key)
	if err != nil {
		return errcode.Annotate(err, "read key")
	}
	endPoint := &creds.Endpoint{
		User:      user,
		Server:    c.Server,
		Key:       keyBytes,
		Transport: c.Transport,
		Homeless:  true,
		NoTTY:     true,
	}

	login, err := creds.NewLogin(endPoint)
	if err != nil {
		return errcode.Annotate(err, "make login")
	}
	token, err := login.Token()
	if err != nil {
		return errcode.Annotate(err, "get token")
	}

	c.TokenSource = httputil.NewStaticToken(token)
	return nil
}
