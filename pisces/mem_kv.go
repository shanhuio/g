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

package pisces

import (
	"sort"
	"sync"

	"shanhu.io/g/errcode"
)

type memKV struct {
	mu sync.RWMutex
	m  map[string]*memEntry
}

func newMemKV() *memKV {
	return &memKV{m: make(map[string]*memEntry)}
}

func (b *memKV) clear() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.m = make(map[string]*memEntry)
	return nil
}

func (b *memKV) add(k, cls string, bs []byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, found := b.m[k]; found {
		return errcode.InvalidArgf("already exist")
	}
	b.m[k] = newMemEntry(cls, bs)
	return nil
}

func (b *memKV) get(k string) ([]byte, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	entry := b.m[k]
	if entry == nil {
		return nil, notFound
	}
	return entry.bytes(), nil
}

func (b *memKV) has(k string) (bool, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	_, ok := b.m[k]
	return ok, nil
}

func (b *memKV) set(k string, bs []byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	entry := b.m[k]
	if entry == nil {
		return notFound
	}
	entry.setBytes(bs)
	return nil
}

func (b *memKV) setClass(k, cls string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	entry := b.m[k]
	if entry == nil {
		return notFound
	}
	entry.cls = cls
	return nil
}

func (b *memKV) remove(k string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.m[k]; !ok {
		return notFound
	}
	delete(b.m, k)
	return nil
}

func (b *memKV) emplace(k, cls string, bs []byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.m[k]; ok {
		return nil
	}
	b.m[k] = newMemEntry(cls, bs)
	return nil
}

func (b *memKV) replace(k, cls string, bs []byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.m[k] = newMemEntry(cls, bs)
	return nil
}

func (b *memKV) appendBytes(k string, bs []byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	entry := b.m[k]
	if entry == nil {
		b.m[k] = newMemEntry("", bs)
	} else {
		entry.appendBytes(bs)
	}
	return nil
}

func (b *memKV) mutate(k string, f func(bs []byte) ([]byte, error)) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	entry := b.m[k]
	if entry == nil {
		return notFound
	}
	bs, err := f(entry.bytes())
	if err != nil {
		return err
	}

	entry.setBytes(bs)
	return nil
}

func (b *memKV) count() (int64, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return int64(len(b.m)), nil
}

func (b *memKV) keys() []string {
	var keys []string
	for k := range b.m {
		keys = append(keys, k)
	}
	return keys
}

func (b *memKV) classKeys(cls string) []string {
	var keys []string
	for k, entry := range b.m {
		if entry.cls == cls {
			keys = append(keys, k)
		}
	}
	return keys
}

func sortKeys(keys []string, desc bool) {
	if !desc {
		sort.Strings(keys)
	} else {
		sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	}
}

func partialKeys(p *KVPartial, keys []string) []string {
	n := uint64(len(keys))
	start := p.Offset
	end := start + p.N
	if start > n {
		start = n
	}
	if end > n {
		end = n
	}
	return keys[start:end]
}

func (b *memKV) walkKeys(keys []string, f WalkFunc) error {
	for _, k := range keys {
		entry := b.m[k]
		if err := f(k, entry.cls, entry.bytes()); err != nil {
			return err
		}
	}
	return nil
}

func (b *memKV) walk(f WalkFunc) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	keys := b.keys()
	sortKeys(keys, false)
	return b.walkKeys(keys, f)
}

func (b *memKV) walkClass(cls string, f WalkFunc) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	keys := b.classKeys(cls)
	sortKeys(keys, false)
	return b.walkKeys(keys, f)
}

func (b *memKV) walkPartial(p *KVPartial, f WalkFunc) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	keys := b.keys()
	sortKeys(keys, p.Desc)
	return b.walkKeys(partialKeys(p, keys), f)
}

func (b *memKV) walkPartialClass(cls string, p *KVPartial, f WalkFunc) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	keys := b.classKeys(cls)
	sortKeys(keys, p.Desc)
	return b.walkKeys(partialKeys(p, keys), f)
}

func (b *memKV) create() error        { return nil }
func (b *memKV) createMissing() error { return nil }
func (b *memKV) destroy() error       { return nil }

func (b *memKV) ops() *KVOps {
	return &KVOps{
		Clear:            b.clear,
		Add:              b.add,
		Get:              b.get,
		Has:              b.has,
		Set:              b.set,
		SetClass:         b.setClass,
		Mutate:           b.mutate,
		Remove:           b.remove,
		Emplace:          b.emplace,
		Replace:          b.replace,
		Append:           b.appendBytes,
		Walk:             b.walk,
		WalkClass:        b.walkClass,
		WalkPartial:      b.walkPartial,
		WalkPartialClass: b.walkPartialClass,
		Count:            b.count,

		Create:        b.create,
		CreateMissing: b.createMissing,
		Destroy:       b.destroy,
	}
}

// NewMemKV creates a new key-value pair in memory.
func NewMemKV() *KV {
	return newKV(newMemKV().ops())
}

// NewOrderedMemKV creates a new ordered key-value pair in memory.
func NewOrderedMemKV() *KV {
	return newOrderedKV(newMemKV().ops())
}
