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

package identity

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"time"

	"shanhu.io/g/errcode"
	"shanhu.io/g/jwt"
	"shanhu.io/g/rsautil"
)

type jwtSigner struct {
	signer Signer
}

func (s *jwtSigner) Header(ctx context.Context) (*jwt.Header, error) {
	id, err := s.signer.Identity(ctx)
	if err != nil {
		return nil, errcode.Annotate(err, "fetch identity")
	}
	if len(id.PublicKeys) == 0 {
		return nil, errcode.Annotate(err, "no signing keys")
	}
	k := id.PublicKeys[len(id.PublicKeys)-1]

	return &jwt.Header{
		Alg:   k.Alg,
		Typ:   jwt.DefaultType,
		KeyID: k.ID,
	}, nil
}

func (s *jwtSigner) Sign(
	ctx context.Context, h *jwt.Header, data []byte,
) ([]byte, error) {
	sig, err := s.signer.Sign(ctx, h.KeyID, data)
	if err != nil {
		return nil, err
	}
	return sig.Sig, nil
}

func publicKeyFromCard(ctx context.Context, card Card, keyID string) (
	*PublicKey, error,
) {
	id, err := card.Identity(ctx)
	if err != nil {
		return nil, errcode.Annotate(err, "fetch identity")
	}
	k := FindPublicKey(id, keyID)
	if k == nil {
		return nil, errcode.NotFoundf("key not found")
	}
	return k, nil
}

type jwtVerifier struct {
	card Card
}

func newJWTVerifier(card Card) *jwtVerifier {
	return &jwtVerifier{card: card}
}

// NewJWTVerifier returns a new JWT verifier using the identity card.
func NewJWTVerifier(card Card) jwt.Verifier {
	return newJWTVerifier(card)
}

func (v *jwtVerifier) Verify(
	ctx context.Context, h *jwt.Header, data, sig []byte, t time.Time,
) error {
	if h.Alg != jwt.AlgRS256 {
		return errcode.InvalidArgf("alg %q not supported", h.Alg)
	}

	k, err := publicKeyFromCard(ctx, v.card, h.KeyID)
	if err != nil {
		return errcode.Annotate(err, "find public key")
	}
	if k.Type != rsaKeyType {
		return errcode.NotFoundf("key type not supported")
	}
	if err := publicKeyValid(k, t); err != nil {
		return errcode.Annotate(err, "invalid key")
	}

	pub, err := rsautil.ParsePublicKey([]byte(k.Key))
	if err != nil {
		return err
	}
	hash := sha256.Sum256(data)
	return rsa.VerifyPKCS1v15(pub, crypto.SHA256, hash[:], sig)
}

func (s *jwtSigner) rsaPublicKeyPEM(ctx context.Context, keyID string) (
	[]byte, error,
) {
	pub, err := publicKeyFromCard(ctx, s.signer, keyID)
	if err != nil {
		return nil, errcode.Annotate(err, "find public key")
	}
	if pub.Type != rsaKeyType {
		return nil, errcode.NotFoundf("key type not supported")
	}
	k, err := rsautil.ParsePublicKey([]byte(pub.Key))
	if err != nil {
		return nil, errcode.Annotate(err, "parse key")
	}
	block := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(k),
	}
	return pem.EncodeToMemory(block), nil
}

func newJWTSigner(signer Signer) *jwtSigner {
	return &jwtSigner{signer: signer}
}

// NewJWTSigner returns a JWT signer with given signer.
func NewJWTSigner(signer Signer) jwt.Signer {
	return newJWTSigner(signer)
}
