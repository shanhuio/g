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
