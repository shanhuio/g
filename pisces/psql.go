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
	"fmt"

	"shanhu.io/pub/sqlx"
)

// PsqlCreateTable creates a postgres table for the given table name and
// scheme.
func PsqlCreateTable(db *sqlx.DB, table, scheme string) error {
	q := fmt.Sprintf(`create table %s %s`, table, scheme)
	_, err := db.X(q)
	return err
}

// PsqlCreateTableMissing creates a postgres table if it does not exist, using
// the given table name and scheme.
func PsqlCreateTableMissing(db *sqlx.DB, table, scheme string) error {
	q := fmt.Sprintf(`create table if not exists %s %s`, table, scheme)
	_, err := db.X(q)
	return err
}

var psqlKVScheme = fmt.Sprintf(`(
	k varchar(%d) primary key not null,
	c varchar(%d) not null,
	v bytea not null
)`, MaxKVKeyLen, MaxKVClassLen)

// PsqlCreateKV creates a key value pair postgres table.
func PsqlCreateKV(db *sqlx.DB, table string) error {
	return PsqlCreateTable(db, table, psqlKVScheme)
}

// PsqlCreateKVMissing creates a key value pair postgres table if the table
// does not exist.
func PsqlCreateKVMissing(db *sqlx.DB, table string) error {
	return PsqlCreateTableMissing(db, table, psqlKVScheme)
}

// PsqlDropExist destroys the table if the table exists. It does nothing if the
// table does not exist.
func PsqlDropExist(db *sqlx.DB, table string) error {
	_, err := db.X(fmt.Sprintf("drop table if exists %s", table))
	return err
}

// PsqlDrop destroys the specific postgres table.
func PsqlDrop(db *sqlx.DB, table string) error {
	_, err := db.X(fmt.Sprintf("drop table %s", table))
	return err
}
