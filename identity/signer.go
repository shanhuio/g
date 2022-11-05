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

// Signature is the result of signing.
type Signature struct {
	KeyID string
	Sig   []byte
}

// Signer provides a read-only interface for signing stuff.
type Signer interface {
	Card

	// Sign signs a blob of data using the given identity key.
	// When key is an empty string, it might use any key to sign.
	Sign(key string, blob []byte) (*Signature, error)
}
