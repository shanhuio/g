package identity

import (
	"context"
)

// Signature is the result of signing.
type Signature struct {
	KeyID string
	Sig   []byte
}

// Signer provides a read-only interface for signing stuff.
type Signer interface {
	Card

	// Sign signs a blob of data using the given identity key.
	// When key is an empty string, it might use any key to sign.
	Sign(ctx context.Context, key string, blob []byte) (*Signature, error)
}
