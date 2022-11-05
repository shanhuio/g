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

import (
	"io"
)

// Objects provides a simple immutable key value store interface
// for saving blob objects.
type Objects interface {
	// Open opens a new blob.
	Open(key string) (io.ReadCloser, error)

	// Create a new object for writing.
	Create(r io.Reader) (string, error)

	// Has checks if an object exists.
	Has(key string) (bool, error)
}
