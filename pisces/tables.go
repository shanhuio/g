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
	"fmt"

	"shanhu.io/pub/errcode"
	"shanhu.io/pub/sqlx"
)

// Tables saves a set of psql tables that are backed by a postgres
// database backend, or backed by memory.
type Tables struct {
	db     *sqlx.DB
	tables []Table
}

// OpenPsqlTables dials into a postgresql connection and creates
// PsqlTables with it.
func OpenPsqlTables(spec string) (*Tables, error) {
	db, err := sqlx.OpenPsql(spec)
	if err != nil {
		return nil, err
	}
	return NewTables(db), nil
}

// OpenSqlite3Tables opens a sqlite3 database using the given file.
func OpenSqlite3Tables(file string) (*Tables, error) {
	db, err := sqlx.OpenSqlite3(file)
	if err != nil {
		return nil, err
	}
	return NewTables(db), nil
}

// NewMemTables creates a new table set that reside entirely inside memory.
func NewMemTables() *Tables {
	return NewTables(nil)
}

// NewTables creates a new table set using the given database backend. When db
// is nil, it uses memory.
func NewTables(db *sqlx.DB) *Tables {
	return &Tables{db: db}
}

// DB returns the underlying database link.
func (ts *Tables) DB() *sqlx.DB { return ts.db }

// Add adds a table into the table set.
func (ts *Tables) Add(t Table) { ts.tables = append(ts.tables, t) }

func (ts *Tables) newOrderedKV(table string) *KV {
	switch driver := tableDriver(ts.db); driver {
	case "":
		return NewOrderedMemKV()
	case sqlx.Psql:
		return NewOrderedPsqlKV(ts.db, table)
	case sqlx.Sqlite3, sqlx.SqliteGo:
		return NewOrderedSqlite3KV(ts.db, table)
	default:
		panic(fmt.Sprintf("unknown database driver: %q", driver))
	}
}

func (ts *Tables) newKV(table string) *KV {
	switch driver := tableDriver(ts.db); driver {
	case "":
		return NewMemKV()
	case sqlx.Psql:
		return NewPsqlKV(ts.db, table)
	case sqlx.Sqlite3, sqlx.SqliteGo:
		return NewSqlite3KV(ts.db, table)
	default:
		panic(fmt.Sprintf("unknown database driver: %q", driver))
	}
}

// NewKV creates a key-value pair table.
func (ts *Tables) NewKV(table string) *KV {
	kv := ts.newKV(table)
	ts.Add(kv)
	return kv
}

// NewOrderedKV creates an ordered key-value pair table.
func (ts *Tables) NewOrderedKV(table string) *KV {
	kv := ts.newOrderedKV(table)
	ts.Add(kv)
	return kv
}

func (ts *Tables) runOnAll(f func(t Table) error) error {
	if len(ts.tables) == 0 {
		return errcode.Internalf("no table")
	}
	for _, t := range ts.tables {
		if err := f(t); err != nil {
			return err
		}
	}
	return nil
}

// Create creates all tables.
func (ts *Tables) Create() error {
	return ts.runOnAll(func(t Table) error { return t.Create() })
}

// CreateMissing creates all missing tables.
func (ts *Tables) CreateMissing() error {
	return ts.runOnAll(func(t Table) error { return t.CreateMissing() })
}

// Destroy drops all tables.
func (ts *Tables) Destroy() error {
	return ts.runOnAll(func(t Table) error { return t.Destroy() })
}
