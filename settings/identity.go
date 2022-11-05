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
