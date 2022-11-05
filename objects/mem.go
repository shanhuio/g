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
	"sync"

	"shanhu.io/pub/hashutil"
)

type mem struct {
	blobs map[string][]byte
	mu    sync.Mutex
}

func newMem() *mem {
	return &mem{blobs: make(map[string][]byte)}
}

// NewMemStore creates a new in-memory object store.
func NewMemStore() Store { return newMem() }

// NewMem creates a new in-memory object store that support streaming.
func NewMem() Objects { return newMem() }

func (m *mem) Open(h string) (io.ReadCloser, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	bs, found := m.blobs[h]
	if !found {
		return nil, notFound(h)
	}

	return ioutil.NopCloser(bytes.NewBuffer(bs)), nil
}

func (m *mem) Create(r io.Reader) (string, error) {
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	return m.put(bs)
}

func (m *mem) Put(bs []byte) (string, error) {
	cp := make([]byte, len(bs))
	copy(cp, bs)
	return m.put(cp)
}

func (m *mem) put(bs []byte) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	h := hashutil.Hash(bs)
	m.blobs[h] = bs
	return h, nil
}

func (m *mem) Get(h string) ([]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	bs, found := m.blobs[h]
	if !found {
		return nil, notFound(h)
	}
	return bs, nil
}

func (m *mem) Has(h string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, found := m.blobs[h]
	return found, nil
}
