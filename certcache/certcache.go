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

package certcache

import (
	"context"

	"golang.org/x/crypto/acme/autocert"
	"shanhu.io/pub/errcode"
	"shanhu.io/pub/pisces"
)

// CertCache is a certificate cache, backed by a database table.
type CertCache struct {
	t *pisces.KV
}

type certEntry struct {
	Content string // We use string as certs are mostly readable characters.
}

// New creates a new certificate cache. When name is empty, "cert_cache" is
// used.
func New(b *pisces.Tables, name string) *CertCache {
	if name == "" {
		name = "cert_cache"
	}
	return &CertCache{t: b.NewKV(name)}
}

// Get gets a certificate entry by name. Returns autocert.ErrCacheMiss if the
// entry is not found.
func (b *CertCache) Get(_ context.Context, name string) ([]byte, error) {
	entry := new(certEntry)
	if err := b.t.Get(name, entry); err != nil {
		if errcode.IsNotFound(err) {
			return nil, autocert.ErrCacheMiss
		}
		return nil, err
	}
	return []byte(entry.Content), nil
}

// Put saves a certificate entry by name. If the entry already exists,
// it will be replaces.
func (b *CertCache) Put(_ context.Context, name string, data []byte) error {
	entry := &certEntry{Content: string(data)}
	return b.t.Replace(name, entry)
}

// Delete remotes an entry. If the entry does not exist, it is a noop.
func (b *CertCache) Delete(_ context.Context, name string) error {
	if err := b.t.Remove(name); err != nil {
		if errcode.IsNotFound(err) {
			return nil
		}
		return err
	}
	return nil
}
