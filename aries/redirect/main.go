package redirect

import (
	"flag"
	"log"

	"shanhu.io/g/aries"
)

type config struct {
	RedirectToDomain string
}

type server struct {
	c *config
}

func (s *server) redirect(c *aries.C) error {
	u := *c.Req.URL // make a shallow copy
	u.Scheme = "https"
	u.Host = s.c.RedirectToDomain
	c.Redirect(u.String())
	return nil
}

func newServer(c *config) (aries.Func, error) {
	s := &server{c: c}
	return s.redirect, nil
}

// Main is the main entrance for the redirect service.
func Main() {
	addr := aries.DeclareAddrFlag("localhost:8000")
	to := flag.String("to", "", "redirect to this address")
	flag.Parse()

	config := &config{RedirectToDomain: *to}

	f, err := newServer(config)
	if err != nil {
		log.Fatal(err)
	}
	if err := aries.ListenAndServe(*addr, f); err != nil {
		log.Fatal(err)
	}
}
