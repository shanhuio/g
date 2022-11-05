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
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"shanhu.io/pub/errcode"
	"shanhu.io/pub/hashutil"
)

func isValidKey(k string) bool {
	if len(k) != 64 {
		return false
	}

	for _, r := range k {
		if r >= 'a' && r <= 'z' {
			continue
		}
		if r >= '0' && r <= '9' {
			continue
		}
		return false
	}

	return true
}

// fsObjects is an object store backed by the file system.
type fsObjects struct {
	tmpDir string
	dir    string
	mu     sync.RWMutex
}

func newFSObjects(dir string) (*fsObjects, error) {
	tmpDir := filepath.Join(dir, "tmp")
	if err := os.MkdirAll(tmpDir, 0700); err != nil {
		return nil, err
	}

	return &fsObjects{
		tmpDir: tmpDir,
		dir:    dir,
	}, nil
}

// NewFS returns a new file system based object store.
func NewFS(dir string) (Objects, error) {
	ret, err := newFSObjects(dir)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (b *fsObjects) filename(key string) string {
	return filepath.Join(b.dir, key)
}

func (b *fsObjects) open(key string) (*os.File, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return os.Open(b.filename(key))
}

func (b *fsObjects) Open(key string) (io.ReadCloser, error) {
	if !isValidKey(key) {
		return nil, errcode.NotFoundf("%q is not a valid key", key)
	}

	f, err := b.open(key)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, notFound(key)
		}
		return nil, err
	}
	return f, nil
}

func (b *fsObjects) commit(k string, f *os.File) error {
	target := b.filename(k)

	b.mu.Lock()
	defer b.mu.Unlock()

	has, err := hasFile(target)
	if err != nil {
		return err
	}
	if has {
		return os.Remove(f.Name())
	}
	return os.Rename(f.Name(), target)
}

func (b *fsObjects) Create(r io.Reader) (string, error) {
	f, err := createTemp(b.tmpDir)
	if err != nil {
		return "", err
	}

	defer func() {
		if f != nil {
			f.Close()
			os.Remove(f.Name())
		}
	}()

	tee := io.TeeReader(r, f)
	k, err := hashutil.HashReader(tee)
	if err != nil {
		return "", err
	}
	if err := f.Close(); err != nil {
		return "", err
	}
	if !isValidKey(k) {
		panic(fmt.Sprintf("invalid key generated: %s", k))
	}

	if err := b.commit(k, f); err != nil {
		return "", err
	}

	f = nil
	return k, nil
}

func hasFile(filename string) (bool, error) {
	s, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return s.Mode().IsRegular(), nil
}

func (b *fsObjects) Has(key string) (bool, error) {
	if !isValidKey(key) {
		return false, nil
	}

	b.mu.RLock()
	defer b.mu.RUnlock()

	return hasFile(b.filename(key))
}
