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

package secrets

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"shanhu.io/pub/errcode"
)

// CheckKey checks if the secret key is a valid one.
func CheckKey(k string) error {
	if k == "" {
		return errcode.InvalidArgf("empty secret key")
	}
	if k[0] == '.' {
		return errcode.InvalidArgf("secret key start with dot: %q", k)
	}
	for _, r := range k {
		if r >= '0' && r <= '9' {
			continue
		}
		if r >= 'a' && r <= 'z' {
			continue
		}
		if r >= 'A' && r <= 'Z' {
			continue
		}
		if r == '_' || r == '-' || r == '.' {
			continue
		}
		return errcode.InvalidArgf("invalid secret key: %q", k)
	}
	return nil
}

// Secrets is a secret store
type Secrets interface {
	Get(k string) ([]byte, error)
}

type memSecrets struct {
	m map[string]string
}

func (s *memSecrets) Get(k string) ([]byte, error) {
	if err := CheckKey(k); err != nil {
		return nil, err
	}
	v, ok := s.m[k]
	if !ok {
		return nil, errcode.NotFoundf("secret %q not found", k)
	}
	return []byte(v), nil
}

// NewMem creates a new memory based secret store.
func NewMem(m map[string]string) Secrets {
	return &memSecrets{m: m}
}

type dirSecrets struct {
	dir string
}

func (s *dirSecrets) Get(k string) ([]byte, error) {
	if err := CheckKey(k); err != nil {
		return nil, err
	}
	f := filepath.Join(s.dir, k)
	bs, err := os.ReadFile(f)
	if err != nil {
		return nil, errcode.FromOS(err)
	}
	return bs, nil
}

// NewDir creates a secret store based on a directory
func NewDir(dir string) Secrets {
	return &dirSecrets{dir: dir}
}

// JSON reads a secret and JSON unmarshal it into v using JSON encoding.
func JSON(s Secrets, k string, v interface{}) error {
	bs, err := s.Get(k)
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, v)
}

// Token reads a secret as a token string, where white spaces are trimmed.
func Token(s Secrets, k string) (string, error) {
	bs, err := s.Get(k)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bs)), nil
}

// TokenDefault reads a secret as a token string, similar to Token(), but
// returns def when the secret is not found.
func TokenDefault(s Secrets, k, def string) (string, error) {
	tok, err := Token(s, k)
	if errcode.IsNotFound(err) {
		return def, nil
	}
	return tok, err
}
