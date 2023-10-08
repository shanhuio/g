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

package argon2

import (
	"crypto/subtle"
	"io"

	"golang.org/x/crypto/argon2"
	"shanhu.io/g/errcode"
)

// Password saves the password hashed with Argon2 algorithm.
type Password struct {
	Key     []byte // Hashed key.
	Salt    []byte
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
}

func readSalt(r io.Reader) ([]byte, error) {
	const saltLen = 16
	salt := make([]byte, saltLen)
	if _, err := r.Read(salt); err != nil {
		return nil, err
	}
	return salt, nil
}

// NewPassword creates a new password encrypted by argon2 algorithm.
func NewPassword(password []byte, rand io.Reader) (*Password, error) {
	// Reference: https://tools.ietf.org/id/draft-irtf-cfrg-argon2-05.html
	salt, err := readSalt(rand)
	if err != nil {
		return nil, errcode.Annotate(err, "create salt")
	}

	const (
		time    = 3
		mem     = 32 * 1024
		threads = 4
		keyLen  = 32
	)

	k := argon2.Key(password, salt, time, mem, threads, keyLen)
	return &Password{
		Key:     k,
		Salt:    salt,
		Time:    time,
		Memory:  mem,
		Threads: threads,
		KeyLen:  keyLen,
	}, nil
}

// Check checks if the password matches.
func (p *Password) Check(password []byte) bool {
	k := argon2.Key(password, p.Salt, p.Time, p.Memory, p.Threads, p.KeyLen)
	return subtle.ConstantTimeCompare(k, p.Key) == 1
}

// CheckString checks if the password matches the given string.
func (p *Password) CheckString(password string) bool {
	return p.Check([]byte(password))
}
