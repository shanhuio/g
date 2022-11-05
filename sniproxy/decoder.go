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
	"fmt"
	"io"
	"io/ioutil"
)

type decoder struct {
	r    io.Reader
	n    int64
	err  error
	tail int64
}

func newDecoder(r io.Reader) *decoder {
	return &decoder{r: r}
}

func (d *decoder) hasErr() bool { return d.err != nil }

func (d *decoder) Err() error { return d.err }

func (d *decoder) count() int64 { return d.n }

func (d *decoder) overread() bool {
	return d.err == io.ErrUnexpectedEOF
}

func (d *decoder) read(buf []byte) {
	if d.err != nil {
		return
	}
	n, err := io.ReadFull(d.r, buf)
	d.n += int64(n)
	if err == io.EOF {
		d.err = io.ErrUnexpectedEOF
	} else if err != nil {
		d.err = err
	}
}

func (d *decoder) rest() []byte {
	if d.err != nil {
		return nil
	}
	bs, err := ioutil.ReadAll(d.r)
	d.n += int64(len(bs))
	d.err = err
	return bs
}

func (d *decoder) u8() byte {
	var buf [1]byte
	d.read(buf[:])
	return buf[0]
}

func (d *decoder) u64() uint64 {
	if d.hasErr() {
		return 0 // To shortcut the encoding part.
	}
	var buf [8]byte
	d.read(buf[:])
	v := endian.Uint64(buf[:])
	return v
}

func (d *decoder) bytes(buf []byte) []byte {
	n := int(d.u64())
	if n <= 0 {
		return nil
	}
	if d.hasErr() {
		return nil // To avoid the allocation.
	}
	if len(buf) >= n {
		buf = buf[:n]
	} else {
		buf = make([]byte, n)
	}
	d.read(buf)
	return buf
}

func (d *decoder) str() string {
	return string(d.bytes(nil))
}

func (d *decoder) tailError() error {
	if d.tail == 0 {
		return nil
	}
	return &tailError{n: d.tail}
}

func (d *decoder) end() {
	if d.hasErr() {
		return
	}

	var buf1 [1]byte
	buf := buf1[:]
	for {
		n, err := d.r.Read(buf)
		d.tail += int64(n)
		if err == io.EOF {
			break
		}
		if err != nil {
			d.err = err
			return
		}
		if len(buf) <= 1 {
			buf = make([]byte, 1024)
		}
	}

	d.err = d.tailError()
	return
}

type tailError struct {
	n int64
}

func (e *tailError) Error() string {
	return fmt.Sprintf("unexpect tail of %d bytes", e.n)
}
