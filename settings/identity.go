package settings

// Identity implements an identity store that
type Identity struct {
	settings Settings
	key      string
}

// NewIdentity create a new identity store.
func NewIdentity(s Settings, k string) *Identity {
	if k == "" {
		k = "identity"
	}
	return &Identity{settings: s, key: k}
}

// Load loads a value.
func (s *Identity) Load(v interface{}) error {
	return s.settings.Get(s.key, v)
}

// Check checks if a value is set already.
func (s *Identity) Check() (bool, error) {
	return s.settings.Has(s.key)
}

// Save saves the value into the store.
func (s *Identity) Save(v interface{}) error {
	return s.settings.Set(s.key, v)
}
