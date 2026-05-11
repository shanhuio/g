package creds

import (
	"crypto/rsa"

	"shanhu.io/g/rsautil"
)

func parsePrivateKey(name string, bs []byte, tty bool) (
	*rsa.PrivateKey, error,
) {
	if tty {
		return rsautil.ParsePrivateKeyTTY(name, bs)
	}
	return rsautil.ParsePrivateKey(bs)
}

func readPrivateKey(pemFile string, tty bool) (*rsa.PrivateKey, error) {
	if tty {
		return rsautil.ReadPrivateKeyTTY(pemFile)
	}
	return rsautil.ReadPrivateKey(pemFile)
}

func readEndpointKey(p *Endpoint) (*rsa.PrivateKey, error) {
	tty := !p.NoTTY
	if p.Key != nil {
		return parsePrivateKey("key", p.Key, tty)
	}
	return readPrivateKey(p.PemFile, tty)
}
