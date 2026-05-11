package rsautil

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"golang.org/x/crypto/ssh"
)

func pemBlock(k *rsa.PrivateKey, pwd []byte) (*pem.Block, error) {
	const pemType = "RSA PRIVATE KEY"

	if pwd == nil {
		return &pem.Block{
			Type:  pemType,
			Bytes: x509.MarshalPKCS1PrivateKey(k),
		}, nil
	}

	return x509.EncryptPEMBlock(
		rand.Reader, pemType,
		x509.MarshalPKCS1PrivateKey(k),
		pwd, x509.PEMCipherDES,
	)
}

// GenerateKey generates a private/public key pair with the given passphrase.
// n is the bit size of the RSA key. When n is less than 0, 4096 is used.
func GenerateKey(passphrase []byte, n int) (pri, pub []byte, err error) {
	if n <= 0 {
		n = 4096
	}
	key, err := rsa.GenerateKey(rand.Reader, n)
	if err != nil {
		return nil, nil, err
	}

	b, err := pemBlock(key, passphrase)
	if err != nil {
		return nil, nil, err
	}

	priBuf := new(bytes.Buffer)
	if err := pem.Encode(priBuf, b); err != nil {
		return nil, nil, err
	}
	pubKey, err := ssh.NewPublicKey(&key.PublicKey)
	if err != nil {
		return nil, nil, err
	}
	pub = ssh.MarshalAuthorizedKey(pubKey)
	return priBuf.Bytes(), pub, nil
}
