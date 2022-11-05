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
	"encoding/json"
)

// MaxKVClassLen is the maximum length of the class string of a hashed KV.
const MaxKVClassLen = 255

// MaxKVKeyLen is the maximum key length of ordered key-value pair storage.
const MaxKVKeyLen = 255

// KV provides a key-value pair table.
type KV struct {
	ops     *KVOps
	ordered bool
}

func newKV(ops *KVOps) *KV {
	return &KV{ops: ops}
}

func newOrderedKV(ops *KVOps) *KV {
	return &KV{ops: ops, ordered: true}
}

func (b *KV) mapKey(k string) (string, error) {
	return kvMapKey(k, b.ordered)
}

// Create creates the table.
func (b *KV) Create() error { return b.ops.Create() }

// CreateMissing creates the table if the table is missing.
func (b *KV) CreateMissing() error { return b.ops.CreateMissing() }

// Destroy destroys the table.
func (b *KV) Destroy() error { return b.ops.Destroy() }

// Clear clears the entire table.
func (b *KV) Clear() error { return b.ops.Clear() }

// AddClass adds an entry with a particular class.
func (b *KV) AddClass(k, cls string, v interface{}) error {
	bs, err := json.Marshal(v)
	if err != nil {
		return err
	}
	mk, err := b.mapKey(k)
	if err != nil {
		return err
	}
	return b.ops.Add(mk, cls, bs)
}

// SetClass set an entry's class string.
func (b *KV) SetClass(k, cls string) error {
	mk, err := b.mapKey(k)
	if err != nil {
		return err
	}
	return b.ops.SetClass(mk, cls)
}

// Add is a short-hand for AddClass but with cls set to empty string.
func (b *KV) Add(k string, v interface{}) error {
	return b.AddClass(k, "", v)
}

// Remove removes the entry with the specific key.
func (b *KV) Remove(k string) error {
	mk, err := b.mapKey(k)
	if err != nil {
		return err
	}
	return b.ops.Remove(mk)
}

// GetBytes gets the value bytes for the specific key.
func (b *KV) GetBytes(k string) ([]byte, error) {
	mk, err := b.mapKey(k)
	if err != nil {
		return nil, err
	}
	return b.ops.Get(mk)
}

// Get gets the value and JSON marshals it into v.
func (b *KV) Get(k string, v interface{}) error {
	bs, err := b.GetBytes(k)
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, v)
}

// Has checks if there exists an entry with the particular key.
func (b *KV) Has(k string) (bool, error) {
	mk, err := b.mapKey(k)
	if err != nil {
		return false, err
	}
	return b.ops.Has(mk)
}

// Emplace sets the value for a particular key.
// Does nothing if the key already exists.
func (b *KV) Emplace(k string, v interface{}) error {
	mk, err := b.mapKey(k)
	if err != nil {
		return err
	}
	bs, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return b.ops.Emplace(mk, "", bs)
}

// Replace sets the value for a particular key. Creates the key if not
// exist.
func (b *KV) Replace(k string, v interface{}) error {
	mk, err := b.mapKey(k)
	if err != nil {
		return err
	}
	bs, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return b.ops.Replace(mk, "", bs)
}

// AppendBytes appends the byte slice to the existing value of the entry
// of the specified key. Creates the key if not exist.
func (b *KV) AppendBytes(k string, bs []byte) error {
	mk, err := b.mapKey(k)
	if err != nil {
		return err
	}
	return b.ops.Append(mk, bs)
}

// SetBytes updates the value bytes of a particular entry.
func (b *KV) SetBytes(k string, bs []byte) error {
	mk, err := b.mapKey(k)
	if err != nil {
		return err
	}
	return b.ops.Set(mk, bs)
}

// Set updates the JSON value of a particular entry.
func (b *KV) Set(k string, v interface{}) error {
	bs, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return b.SetBytes(k, bs)
}

// Mutate applies a function to an entry's value.
func (b *KV) Mutate(
	k string, v interface{}, f func(v interface{}) error,
) error {
	mk, err := b.mapKey(k)
	if err != nil {
		return err
	}

	err = b.ops.Mutate(mk, func(bs []byte) ([]byte, error) {
		if err := json.Unmarshal(bs, v); err != nil {
			return nil, err
		}
		if err := f(v); err != nil {
			return nil, err
		}
		return json.Marshal(v)
	})
	if err == ErrCancel {
		return nil
	}
	return err
}

// Count returns the total number of entries.
func (b *KV) Count() (int64, error) { return b.ops.Count() }

// Walk iterates through all entriess in the key-value store.
func (b *KV) Walk(it *Iter) error {
	err := b.ops.Walk(it.doWalk)
	if err == ErrCancel {
		return nil
	}
	return err
}

// WalkClass iterates through all entries in the key-value store of
// a particular class.
func (b *KV) WalkClass(cls string, it *Iter) error {
	err := b.ops.WalkClass(cls, it.doWalk)
	if err == ErrCancel {
		return nil
	}
	return err
}

// WalkPartial iterates through a part of the entries in the key-value store,
// specificed by the partial option.
func (b *KV) WalkPartial(p *KVPartial, it *Iter) error {
	if !b.ordered {
		return ErrUnordered
	}
	err := b.ops.WalkPartial(p, it.doWalk)
	if err == ErrCancel {
		return nil
	}
	return err
}

// WalkPartialClass iterates through a part of the entries of a particular
// class in the key-value store, specified by the partial option.
func (b *KV) WalkPartialClass(cls string, p *KVPartial, it *Iter) error {
	if !b.ordered {
		return ErrUnordered
	}
	err := b.ops.WalkPartialClass(cls, p, it.doWalk)
	if err == ErrCancel {
		return nil
	}
	return err
}
