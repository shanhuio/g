// Copyright (C) 2022  Shanhu Tech Inc.
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
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"time"

	"shanhu.io/pub/errcode"
	"shanhu.io/pub/hashutil"
	"shanhu.io/pub/jwt"
	"shanhu.io/pub/rsautil"
	"shanhu.io/pub/timeutil"
)

type simpleData struct {
	Identity    *Identity     // Public identity data.
	PrivateKeys []*privateKey `json:",omitempty"`
}

type privateKey struct {
	ID  string // Related key id.
	Key string // The private part.
}

// SimpleStore is a simple store for saving / loading data.
// This can be used to implement a simple identity core.
type SimpleStore interface {
	// Check checks if some data has saved yet.
	Check() (bool, error)

	// Save saves the data.
	Save(v interface{}) error

	// Load loads the data.
	Load(v interface{}) error
}

type simpleCore struct {
	store SimpleStore
	now   func() time.Time
}

// NewSimpleCore creates a new simple core using the given store.
// simple store only support RS256 signing.
func NewSimpleCore(store SimpleStore, t func() time.Time) Core {
	return &simpleCore{
		store: store,
		now:   timeutil.NowFunc(t),
	}
}

const rsaKeyType = "ssh-rsa"

func (c *simpleCore) Init(config *CoreConfig) (*Identity, error) {
	check, err := c.store.Check()
	if err != nil {
		return nil, errcode.Annotate(err, "check key")
	}
	if check {
		return nil, ErrAlreadyInitialized
	}

	if len(config.Keys) == 0 {
		return nil, errcode.InvalidArgf("must init at least one key")
	}

	now := c.now()
	for i, k := range config.Keys {
		if k.Type != "" {
			return nil, errcode.InvalidArgf("key #%d: type not supported", i)
		}
		if k.NotValidAfter == 0 {
			return nil, errcode.InvalidArgf("key #%d: missing expire time", i)
		}
		expire := time.Unix(k.NotValidAfter, 0)
		if expire.Before(now) {
			return nil, errcode.InvalidArgf("key #%d: already expired", i)
		}
		if k.NotValidBefore != 0 && k.NotValidBefore >= k.NotValidAfter {
			return nil, errcode.InvalidArgf("key #%d: never valid", i)
		}
	}

	id := new(Identity)
	var privateKeys []*privateKey
	const keySize = 2048
	for i, k := range config.Keys {
		pri, pub, err := rsautil.GenerateKey(nil, keySize)
		if err != nil {
			return nil, errcode.Internalf("generate key #%d", i)
		}

		keyID := hashutil.Hash(pub)
		id.PublicKeys = append(id.PublicKeys, &PublicKey{
			ID:             keyID,
			Type:           rsaKeyType,
			Alg:            jwt.AlgRS256,
			Key:            string(pub),
			NotValidAfter:  k.NotValidAfter,
			NotValidBefore: k.NotValidBefore,
			Comment:        k.Comment,
		})
		privateKeys = append(privateKeys, &privateKey{
			ID:  keyID,
			Key: string(pri),
		})
	}

	data := &simpleData{
		Identity:    id,
		PrivateKeys: privateKeys,
	}

	if err := c.store.Save(data); err != nil {
		return nil, errcode.Annotate(err, "save new identity")
	}
	return id, nil
}

func (c *simpleCore) Identity(ctx context.Context) (*Identity, error) {
	id := new(simpleData)
	if err := c.store.Load(id); err != nil {
		return nil, errcode.Annotate(err, "load identity")
	}
	// Clear the private keys, so that they can be garbage collected.
	id.PrivateKeys = nil
	return id.Identity, nil
}

func (c *simpleCore) AddKey(config *KeyConfig) (*PublicKey, error) {
	return nil, errcode.Internalf("not implemented")
}

func (c *simpleCore) RemoveKey(id string) error {
	dat := new(simpleData)
	if err := c.store.Load(dat); err != nil {
		return errcode.Annotate(err, "load identity")
	}
	if dat.Identity == nil {
		return errcode.Internalf("identity missing")
	}

	pubKeyMap := make(map[string]*PublicKey)
	for _, k := range dat.Identity.PublicKeys {
		pubKeyMap[k.ID] = k
	}

	var (
		privKeys []*privateKey
		pubKeys  []*PublicKey
	)
	for _, k := range dat.PrivateKeys {
		if k.ID == id {
			continue
		}
		pub, ok := pubKeyMap[k.ID]
		if !ok {
			continue
		}
		privKeys = append(privKeys, k)
		pubKeys = append(pubKeys, pub)
	}

	if len(privKeys) == 0 {
		return errcode.InvalidArgf("no key left after removal")
	}

	dat.PrivateKeys = privKeys
	dat.Identity.PublicKeys = pubKeys

	if err := c.store.Save(dat); err != nil {
		return errcode.Annotate(err, "save identity")
	}
	return nil
}

func (c *simpleCore) Sign(ctx context.Context, key string, blob []byte) (
	*Signature, error,
) {
	id := new(simpleData)
	if err := c.store.Load(id); err != nil {
		return nil, errcode.Annotate(err, "load identity")
	}

	// Sanity check on data.
	if id.Identity == nil {
		return nil, errcode.Internalf("identity missing")
	}
	if len(id.PrivateKeys) == 0 {
		return nil, errcode.Internalf("no key to sign")
	}

	// Pick the private key.
	var pri *privateKey
	if key == "" {
		// When key id not specified, always use the last key.
		pri = id.PrivateKeys[len(id.PrivateKeys)-1]
		key = pri.ID
	} else {
		// Pick private key based on ID.
		for _, k := range id.PrivateKeys {
			if k.ID == key {
				pri = k
				break
			}
		}
		if pri == nil {
			return nil, errcode.NotFoundf("key not found")
		}
	}

	// Find the corresponding public key.
	var pub *PublicKey
	for _, k := range id.Identity.PublicKeys {
		if k.ID == key {
			pub = k
			break
		}
	}
	if pub == nil {
		return nil, errcode.Internalf("public key not found")
	}

	// Check key type and validity.
	if pub.Type != rsaKeyType {
		return nil, errcode.Internalf("unknown key type: %s", pub.Type)
	}

	now := c.now()
	if err := publicKeyValid(pub, now); err != nil {
		return nil, errcode.Annotate(err, "invalid key")
	}

	// Parse the private key.
	rsaKey, err := rsautil.ParsePrivateKey([]byte(pri.Key))
	if err != nil {
		return nil, errcode.Annotate(err, "parse key")
	}

	// Sign!
	hash := sha256.Sum256(blob)
	sig, err := rsa.SignPKCS1v15(
		rand.Reader, rsaKey, crypto.SHA256, hash[:],
	)
	if err != nil {
		return nil, err
	}

	return &Signature{KeyID: key, Sig: sig}, nil
}
