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

package settings

import (
	"shanhu.io/pub/pisces"
)

// Table is a pisces.KV based settings implementation.
type Table struct {
	t *pisces.KV
}

// NewTableName creates a new settings table using the given table name.
func NewTableName(b *pisces.Tables, name string) *Table {
	return &Table{t: b.NewKV(name)}
}

// NewTable creates a new settings table that is named as "settings".
func NewTable(b *pisces.Tables) *Table {
	return NewTableName(b, "settings")
}

// Get gets a settings.
func (b *Table) Get(key string, v interface{}) error {
	return b.t.Get(key, v)
}

// Has checks if a setting exists.
func (b *Table) Has(key string) (bool, error) {
	return b.t.Has(key)
}

// Set sets the value of a settings key.
func (b *Table) Set(key string, v interface{}) error {
	return b.t.Replace(key, v)
}
