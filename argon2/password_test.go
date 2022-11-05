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

package argon2

import (
	"testing"

	"crypto/rand"
	"errors"
	"io"
)

func TestPassword(t *testing.T) {
	pass := []byte("my password")

	ar, err := NewPassword(pass, rand.Reader)
	if err != nil {
		t.Fatal("hash with argon2: ", err)
	}

	if !ar.Check(pass) {
		t.Error("check password failed")
	}

	if !ar.CheckString(string(pass)) {
		t.Error("check password string failed")
	}

	if ar.Check([]byte("wrong password")) {
		t.Error("check wrong password passed")
	}
}

type errorReader struct{}

func (r *errorReader) Read([]byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

func TestPassword_badRand(t *testing.T) {
	r := new(errorReader)

	ar, err := NewPassword([]byte("pwd"), r)
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Errorf("got %s, want %s", err, io.ErrUnexpectedEOF)
		t.Logf("password %x", ar)
	}
}
