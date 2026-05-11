package oauth2

import (
	"shanhu.io/g/aries"
)

type signInHandler struct {
	client   *Client
	redirect string
}

func newSignInHandler(client *Client, redirect string) *signInHandler {
	return &signInHandler{
		client:   client,
		redirect: redirect,
	}
}

func (h *signInHandler) Serve(c *aries.C) error {
	redirect := h.redirect
	if r := c.Req.URL.Query().Get("r"); r != "" {
		parsed, err := ParseRedirect(r)
		if err != nil {
			return err
		}
		redirect = parsed
	}
	state := &State{Dest: redirect}
	c.Redirect(h.client.SignInURL(state))
	return nil
}
