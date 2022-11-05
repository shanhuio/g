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
	"errors"
	"time"
)

// KeyConfig is the configuration for a new key.
type KeyConfig struct {
	Type           string // Optional
	NotValidAfter  int64
	NotValidBefore int64
	Comment        string
}

// CoreConfig is the configuration for initialiazation of the identity.
type CoreConfig struct {
	Keys []*KeyConfig
}

// SingleKeyCoreConfig creates a simple CoreConfig that creates one
// single key that expires at expire.
func SingleKeyCoreConfig(expire time.Time) *CoreConfig {
	return &CoreConfig{
		Keys: []*KeyConfig{{NotValidAfter: expire.Unix()}},
	}
}

// Core is an identity core that can save the identity keys.
type Core interface {
	// Init initializes the identity with the given config.
	Init(c *CoreConfig) (*Identity, error)

	// AddKey adds a new identity key.
	AddKey(c *KeyConfig) (*PublicKey, error)

	// RemoveKey removes an identity key.
	RemoveKey(id string) error

	Signer
}

// ErrAlreadyInitialized is returned if Init() is called
// when the KeyStore is already initialized.
var ErrAlreadyInitialized = errors.New("already initialized")

// MakeSureInit initializes the core if it is not already initialized.
func MakeSureInit(core Core, config *CoreConfig) error {
	if _, err := core.Init(config); err != nil {
		if err != ErrAlreadyInitialized {
			return err
		}
	}
	return nil
}
