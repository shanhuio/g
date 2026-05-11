package identity

import (
	"shanhu.io/g/aries"
)

type service struct {
	card Card
}

func newService(card Card) *service {
	return &service{card: card}
}

// GetIDRequest is the request for getting an identity.
type GetIDRequest struct{}

func (s *service) apiGet(c *aries.C, req *GetIDRequest) (*Identity, error) {
	return s.card.Identity(c.Context)
}

// NewService creates a new identity service stub
func NewService(card Card) aries.Service {
	s := newService(card)
	r := aries.NewRouter()
	r.Call("get", s.apiGet)
	return r
}
