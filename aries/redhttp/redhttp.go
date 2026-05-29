// Package redhttp provides a service that routes
// all incoming http requests to https.
package redhttp

import (
	"flag"
	"log"

	"shanhu.io/g/aries"
	"shanhu.io/std/errcode"
)

// Redirect redirects the incoming request to https.
func Redirect(c *aries.C) error {
	host := c.Req.Host
	if host == "" {
		return errcode.InvalidArgf("host not found")
	}
	u := *c.Req.URL
	u.Host = host
	u.Scheme = "https"
	c.Redirect(u.String())
	return nil
}

// Main is the main entrance for the service.
func Main() {
	addr := aries.DeclareAddrFlag("localhost:8000")
	flag.Parse()
	if err := aries.ListenAndServe(*addr, aries.Func(Redirect)); err != nil {
		log.Fatal(err)
	}
}
