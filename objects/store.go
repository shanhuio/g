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

package objects

// Store provides a simple immutable key value store interface
// for saving blob objects.
type Store interface {
	// New creates a blob, returning the blob key.
	Put(bs []byte) (string, error)

	// Get fetches the content of a blob, returning the blob content.
	Get(key string) ([]byte, error)

	// Has checks if a blob of the key exists.
	Has(key string) (bool, error)
}
