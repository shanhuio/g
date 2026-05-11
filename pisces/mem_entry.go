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
