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
)

// Identity is the identity of a service or a robot.
type Identity struct {
	PublicKeys []*PublicKey `json:",omitempty"`
}

// Identity returns itself, so it implements the Card interface.
func (id *Identity) Identity(_ context.Context) (*Identity, error) {
	return id, nil
}

// PublicKey is the public key of an identity.
type PublicKey struct {
	ID             string
	Type           string
	Alg            string // Signing alghorithm,must use JWT alg codes.
	Key            string // Key content.
	NotValidAfter  int64
	NotValidBefore int64  `json:",omitempty"`
	Comment        string `json:",omitempty"`
}

// FindPublicKey finds the public key of the given ID.
// Returns nil if not found.
func FindPublicKey(id *Identity, keyID string) *PublicKey {
	var pub *PublicKey
	for _, k := range id.PublicKeys {
		if k.ID == keyID {
			pub = k
			break
		}
	}
	return pub
}

// Card provides the Identity of an entity.
type Card interface {
	// Identity fetches the identity of the service.
	Identity(ctx context.Context) (*Identity, error)
}
