package rsautil

import (
	"bytes"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"

	"golang.org/x/crypto/ssh"
)

func keyHashStr(h []byte) string {
	return base64.RawURLEncoding.EncodeToString(h)
}

// PublicKeyHash returns the public key hash of a key.
func PublicKeyHash(k *rsa.PublicKey) ([]byte, error) {
	sshPub, err := ssh.NewPublicKey(k)
	if err != nil {
		return nil, err
	}

	wire := bytes.TrimSpace(ssh.MarshalAuthorizedKey(sshPub))
	h := sha256.Sum256(wire)
	return h[:], nil
}

// PublicKeyHashString returns the public key hash string of a key.
func PublicKeyHashString(k *rsa.PublicKey) (string, error) {
	h, err := PublicKeyHash(k)
	if err != nil {
		return "", err
	}
	return keyHashStr(h), nil
}
