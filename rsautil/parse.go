// Copyright (C) 2023  Shanhu Tech Inc.
//
// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU Affero General Public License as published by the
// Free Software Foundation, either version 3 of the License, or (at your
// option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.  See the GNU Affero General Public License
// for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package rsautil

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
	"shanhu.io/g/errcode"
	"shanhu.io/g/osutil"
	"shanhu.io/g/termutil"
)

var (
	errNotRSA = errcode.InvalidArgf("public key is not an RSA key")
	errNoKey  = errcode.InvalidArgf("no key")
)

// ParsePrivateKey parses the PEM encoded RSA private key.
func ParsePrivateKey(bs []byte) (*rsa.PrivateKey, error) {
	if len(bs) == 0 {
		return nil, errNoKey
	}

	b, _ := pem.Decode(bs)
	if b == nil {
		return nil, errcode.InvalidArgf("pem decode failed")
	}
	if x509.IsEncryptedPEMBlock(b) {
		return nil, errcode.InvalidArgf("key is encrypted")
	}
	return x509.ParsePKCS1PrivateKey(b.Bytes)
}

// ReadPrivateKey parses the PEM encded RSA private key file.
func ReadPrivateKey(f string) (*rsa.PrivateKey, error) {
	bs, err := osutil.ReadPrivateFile(f)
	if err != nil {
		return nil, err
	}
	return ParsePrivateKey(bs)
}

// ParsePublicKey parses a marshalled public key in SSH authorized key format.
func ParsePublicKey(bs []byte) (*rsa.PublicKey, error) {
	if len(bs) == 0 {
		return nil, errNoKey
	}

	k, _, _, _, err := ssh.ParseAuthorizedKey(bs)
	if err != nil {
		return nil, err
	}

	if k.Type() != "ssh-rsa" {
		return nil, errNotRSA
	}
	ck, ok := k.(ssh.CryptoPublicKey)
	if !ok {
		return nil, errNotRSA
	}

	ret, ok := ck.CryptoPublicKey().(*rsa.PublicKey)
	if !ok {
		return nil, errNotRSA
	}
	return ret, nil
}

// ReadPublicKey parses a marshalled public key file in SSH authorized key
// file format.
func ReadPublicKey(f string) (*rsa.PublicKey, error) {
	bs, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return ParsePublicKey(bs)
}

// ParsePrivateKeyTTY parses a private key and asks for the passphrase
// if the key is an encrypted PEM.
func ParsePrivateKeyTTY(name string, bs []byte) (
	*rsa.PrivateKey, error,
) {
	b, _ := pem.Decode(bs)
	if b == nil {
		return nil, errcode.InvalidArgf("%q decode failed", name)
	}

	if !x509.IsEncryptedPEMBlock(b) {
		return x509.ParsePKCS1PrivateKey(b.Bytes)
	}

	prompt := fmt.Sprintf("Passphrase for %s: ", name)
	pwd, err := termutil.ReadPassword(prompt)
	if err != nil {
		return nil, err
	}
	der, err := x509.DecryptPEMBlock(b, pwd)
	if err != nil {
		return nil, err
	}
	return x509.ParsePKCS1PrivateKey(der)
}

// ReadPrivateKeyTTY reads a private key from a key file.
func ReadPrivateKeyTTY(pemFile string) (*rsa.PrivateKey, error) {
	bs, err := osutil.ReadPrivateFile(pemFile)
	if err != nil {
		return nil, err
	}
	return ParsePrivateKeyTTY(pemFile, bs)
}
