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

package sniproxy

import (
	"io"
)

type encoder struct {
	w   io.Writer
	n   int64
	err error
}

func newEncoder(w io.Writer) *encoder {
	return &encoder{w: w}
}

func (e *encoder) write(bs []byte) {
	if e.hasErr() {
		return
	}
	n, err := e.w.Write(bs)
	if err != nil {
		e.err = err
		return
	}
	e.n += int64(n)
}

func (e *encoder) Err() error { return e.err }

func (e *encoder) hasErr() bool { return e.err != nil }

func (e *encoder) u64(v uint64) {
	if e.hasErr() {
		return // To shortcut the encoding part.
	}
	var bs [8]byte
	endian.PutUint64(bs[:], v)
	e.write(bs[:])
}

func (e *encoder) u8(v uint8) {
	if e.hasErr() {
		return
	}
	e.write([]byte{v})
}

func (e *encoder) bytes(bs []byte) {
	e.u64(uint64(len(bs)))
	e.write(bs)
}

func (e *encoder) str(s string) {
	bs := []byte(s)
	e.bytes(bs)
}
