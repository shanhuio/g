// Copyright (C) 2023  Shanhu Tech Inc.
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

package states

import (
	"net/url"
	"sync"

	"shanhu.io/pub/errcode"
)

type memBack struct {
	mu sync.Mutex
	m  map[string][]byte
}

func newMemBack() *memBack {
	return &memBack{m: make(map[string][]byte)}
}

func copyBytes(bs []byte) []byte {
	cp := make([]byte, len(bs))
	copy(cp, bs)
	return cp
}

func (b *memBack) Get(_ C, key string) ([]byte, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	bs, ok := b.m[key]
	if !ok {
		return nil, errcode.NotFoundf("%q not found", key)
	}
	return copyBytes(bs), nil
}

func (b *memBack) Put(_ C, key string, data []byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.m[key] = copyBytes(data)
	return nil
}

func (b *memBack) Del(_ C, key string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, ok := b.m[key]; !ok {
		return errcode.NotFoundf("%q not found", key)
	}
	delete(b.m, key)
	return nil
}

func (b *memBack) URL() *url.URL {
	return &url.URL{
		Scheme: "memory",
	}
}

// NewMem returns a new states storage backed by memory.
func NewMem() States { return newMemBack() }
