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

// WalkFunc is the a function type for walking through a table query result.
type WalkFunc func(k, cls string, bs []byte) error

// KVOps provides operations to operate over an unordered key-value pair table.
type KVOps struct {
	Clear    func() error
	Add      func(key, cls string, bs []byte) error
	Get      func(key string) ([]byte, error)
	Has      func(key string) (bool, error)
	Set      func(key string, bs []byte) error
	SetClass func(key, cls string) error
	Mutate   func(key string, f func(bs []byte) ([]byte, error)) error
	Remove   func(key string) error
	Emplace  func(key, cls string, bs []byte) error
	Replace  func(key, cls string, bs []byte) error
	Append   func(key string, bs []byte) error

	Walk             func(f WalkFunc) error
	WalkClass        func(cls string, f WalkFunc) error
	WalkPartial      func(p *KVPartial, f WalkFunc) error
	WalkPartialClass func(cls string, p *KVPartial, f WalkFunc) error
	Count            func() (int64, error)

	Create        func() error
	CreateMissing func() error
	Destroy       func() error
}
