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

	"shanhu.io/g/sqlx"
)

// Sqlite3CreateTable creates a table for the given table for the given
// table name and scheme.
func Sqlite3CreateTable(db *sqlx.DB, table, scheme string) error {
	q := fmt.Sprintf(`create table %s %s`, table, scheme)
	_, err := db.X(q)
	return err
}

// Sqlite3CreateTableMissing creates a postgres table if it does not exist,
// using the given table name and scheme.
func Sqlite3CreateTableMissing(db *sqlx.DB, table, scheme string) error {
	q := fmt.Sprintf(`create table if not exists %s %s`, table, scheme)
	_, err := db.X(q)
	return err
}

const sqlite3KVScheme = `(
	k text not null unique,
	c text not null,
	v blob not null
)`

// Sqlite3CreateKV creates a key value pair sqlite3 table.
func Sqlite3CreateKV(db *sqlx.DB, table string) error {
	return Sqlite3CreateTable(db, table, sqlite3KVScheme)
}

// Sqlite3CreateKVMissing creates a key value pair sqlite3 table if the table
// does not exist.
func Sqlite3CreateKVMissing(db *sqlx.DB, table string) error {
	return Sqlite3CreateTableMissing(db, table, sqlite3KVScheme)
}

// Sqlite3DropExist destroys the table if the table exists. It does nothing if
// the table does not exist.
func Sqlite3DropExist(db *sqlx.DB, table string) error {
	_, err := db.X(fmt.Sprintf("drop table if exists %s", table))
	return err
}

// Sqlite3Drop destroys the specific sqlite3 table.
func Sqlite3Drop(db *sqlx.DB, table string) error {
	_, err := db.X(fmt.Sprintf("drop table %s", table))
	return err
}
