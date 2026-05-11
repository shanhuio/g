package settings

import (
	"shanhu.io/g/pisces"
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
