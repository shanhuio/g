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
	"bytes"
	"io"
	"io/ioutil"
)

type mappedStore struct {
	s Store
}

// NewMapped creates a Objects Store from the old-fashioned Store interface.
func NewMapped(s Store) Objects {
	return &mappedStore{s: s}
}

func (b *mappedStore) Open(key string) (io.ReadCloser, error) {
	bs, err := b.s.Get(key)
	if err != nil {
		return nil, err
	}

	return ioutil.NopCloser(bytes.NewReader(bs)), nil
}

func (b *mappedStore) Create(r io.Reader) (string, error) {
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, r); err != nil {
		return "", err
	}

	return b.s.Put(buf.Bytes())
}

func (b *mappedStore) Has(key string) (bool, error) {
	return b.s.Has(key)
}
