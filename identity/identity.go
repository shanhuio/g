package identity

import (
	"context"
)

// Identity is the identity of a service or a robot.
type Identity struct {
	PublicKeys []*PublicKey `json:",omitempty"`
}

// Identity returns itself, so it implements the Card interface.
func (id *Identity) Identity(_ context.Context) (*Identity, error) {
	return id, nil
}

// PublicKey is the public key of an identity.
type PublicKey struct {
	ID             string
	Type           string
	Alg            string // Signing alghorithm,must use JWT alg codes.
	Key            string // Key content.
	NotValidAfter  int64
	NotValidBefore int64  `json:",omitempty"`
	Comment        string `json:",omitempty"`
}

// FindPublicKey finds the public key of the given ID.
// Returns nil if not found.
func FindPublicKey(id *Identity, keyID string) *PublicKey {
	var pub *PublicKey
	for _, k := range id.PublicKeys {
		if k.ID == keyID {
			pub = k
			break
		}
	}
	return pub
}

// Card provides the Identity of an entity.
type Card interface {
	// Identity fetches the identity of the service.
	Identity(ctx context.Context) (*Identity, error)
}
