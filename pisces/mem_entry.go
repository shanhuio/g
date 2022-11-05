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

package pisces

import (
	"bytes"
)

type memEntry struct {
	cls string
	buf *bytes.Buffer
}

func newMemEntry(cls string, bs []byte) *memEntry {
	ret := &memEntry{
		cls: cls,
		buf: new(bytes.Buffer),
	}
	ret.setBytes(bs)
	return ret
}

func (entry *memEntry) setBytes(bs []byte) {
	entry.buf.Truncate(0)
	entry.buf.Write(bs)
}

func (entry *memEntry) bytes() []byte {
	bs := entry.buf.Bytes()
	if len(bs) == 0 {
		return nil
	}
	ret := make([]byte, len(bs))
	copy(ret, bs)
	return ret
}

func (entry *memEntry) appendBytes(bs []byte) {
	entry.buf.Write(bs)
}
