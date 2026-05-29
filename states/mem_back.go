package states

import (
	"net/url"
	"sync"

	"shanhu.io/std/errcode"
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
