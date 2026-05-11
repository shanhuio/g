package rsautil

import (
	"bytes"
	"crypto/rsa"
)

// PublicKey carries a public key.
type PublicKey struct {
	key     *rsa.PublicKey
	hash    []byte
	hashStr string
}

// NewPublicKey parses a new public key from SSH
// authorized key format.
func NewPublicKey(bs []byte) (*PublicKey, error) {
	k, err := ParsePublicKey(bs)
	if err != nil {
		return nil, err
	}
	h, err := PublicKeyHash(k)
	if err != nil {
		return nil, err
	}
	s := keyHashStr(h)

	return &PublicKey{
		key:     k,
		hash:    h,
		hashStr: s,
	}, nil
}

// Key returns the public key parsed from the bytes.
func (k *PublicKey) Key() *rsa.PublicKey { return k.key }

// HashStr returns the base64 encoding of the key hash.
func (k *PublicKey) HashStr() string { return k.hashStr }

// ParsePublicKeys parses a list of public keys.
func ParsePublicKeys(bs []byte) ([]*PublicKey, error) {
	lines := bytes.Split(bs, []byte{'\n'})
	var keys []*PublicKey
	for _, line := range lines {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		k, err := NewPublicKey(line)
		if err != nil {
			return nil, err
		}
		keys = append(keys, k)
	}
	return keys, nil
}
