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

package signin

import (
	"crypto/rsa"
	"encoding/json"
	"time"

	"golang.org/x/crypto/ssh"
	"shanhu.io/pub/aries"
	"shanhu.io/pub/errcode"
	"shanhu.io/pub/rsautil"
	"shanhu.io/pub/signer"
	"shanhu.io/pub/signin/signinapi"
	"shanhu.io/pub/timeutil"
)

// SSHCertExchangeConfig is the configuration to create an SSH certificate
// signin stub.
type SSHCertExchangeConfig struct {
	CAPublicKey     []byte                 `json:",omitempty"`
	CAPublicKeyFunc func() ([]byte, error) `json:",omitempty"`
	CAPublicKeyFile string                 `json:",omitempty"`

	ChallengeKey []byte

	// Time function for checking certificate. It is not used for
	// token generation.
	Now func() time.Time
}

// SSHCertExchange is a service stub that provides session tokens if the
// user signs a challenge and the SSH certificate of it.
type SSHCertExchange struct {
	tokener Tokener

	caPublicKey *rsa.PublicKey
	ch          *Challenger
	nowFunc     func() time.Time
}

func caPublicKeyFromConfig(conf *SSHCertExchangeConfig) (
	*rsa.PublicKey, error,
) {
	if conf.CAPublicKeyFile != "" {
		return rsautil.ReadPublicKey(conf.CAPublicKeyFile)
	}

	var keyBytes = conf.CAPublicKey
	if keyBytes == nil && conf.CAPublicKeyFunc != nil {
		bs, err := conf.CAPublicKeyFunc()
		if err != nil {
			return nil, errcode.Annotate(err, "fetch CA public key")
		}
		keyBytes = bs
	}

	k, err := rsautil.ParsePublicKey(keyBytes)
	if err != nil {
		return nil, errcode.Annotate(err, "parse CA public key")
	}
	return k, nil
}

// NewSSHCertExchange creates a new SSH certificate exchange that exchanges
// signed challenges for session tokens.
func NewSSHCertExchange(tok Tokener, conf *SSHCertExchangeConfig) (
	*SSHCertExchange, error,
) {
	caPubKey, err := caPublicKeyFromConfig(conf)
	if err != nil {
		return nil, errcode.Annotate(err, "read CA public key")
	}
	signer := signer.New(conf.ChallengeKey)
	ch := NewChallenger(&ChallengerConfig{
		Signer: signer,
		Now:    conf.Now,
	})
	return &SSHCertExchange{
		tokener:     tok,
		caPublicKey: caPubKey,
		ch:          ch,
		nowFunc:     timeutil.NowFunc(conf.Now),
	}, nil
}

func (s *SSHCertExchange) apiSignIn(
	c *aries.C, req *signinapi.SSHSignInRequest,
) (*signinapi.Creds, error) {
	record := new(signinapi.SSHSignInRecord)
	if err := json.Unmarshal(req.RecordBytes, record); err != nil {
		return nil, errcode.Annotate(err, "parse signin record")
	}
	user := record.User
	if user == "" {
		return nil, errcode.InvalidArgf("user name is empty")
	}

	if _, err := s.ch.Check(record.Challenge); err != nil {
		return nil, errcode.Annotate(err, "challenge check failed")
	}

	// Parse Certificate.
	certKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(req.Certificate))
	if err != nil {
		return nil, errcode.Annotate(err, "parse certificate")
	}
	cert, ok := certKey.(*ssh.Certificate)
	if !ok {
		return nil, errcode.InvalidArgf("invalid certificate")
	}

	// Check if it is a user certificate.
	if cert.CertType != ssh.UserCert {
		return nil, errcode.InvalidArgf("not a user certificate")
	}

	// Check CA key.
	cryptoPubKey, ok := cert.SignatureKey.(ssh.CryptoPublicKey)
	if !ok {
		return nil, errcode.InvalidArgf("not a crypto public key")
	}
	if !s.caPublicKey.Equal(cryptoPubKey.CryptoPublicKey()) {
		return nil, errcode.Unauthorizedf("unrecognized CA")
	}

	// Check the time and the certificate.
	checker := &ssh.CertChecker{Clock: s.nowFunc}
	if err := checker.CheckCert(user, cert); err != nil {
		return nil, errcode.Annotate(err, "check certificate failed")
	}

	// Check the signature.
	sig := &ssh.Signature{
		Format: req.Sig.Format,
		Blob:   req.Sig.Blob,
		Rest:   req.Sig.Rest,
	}
	if err := cert.Verify(req.RecordBytes, sig); err != nil {
		return nil, errcode.Annotate(err, "check signature")
	}

	// Get a token.
	ttl := timeutil.TimeDuration(record.TTL)
	token := s.tokener.Token(user, ttl)
	return TokenCreds(user, token), nil
}

// API returns the API router stub for signing in with SSH certificate
// credentials.
func (s *SSHCertExchange) API() *aries.Router {
	r := aries.NewRouter()
	r.Call("challenge", s.ch.Serve)
	r.Call("signin", s.apiSignIn)
	return r
}

// AddAPI adds the API to under /ssh .
func (s *SSHCertExchange) AddAPI(r *aries.Router) {
	r.DirService("ssh", s.API())
}
