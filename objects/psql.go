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

package objects

import (
	"shanhu.io/g/hashutil"
	"shanhu.io/g/pathutil"
	"shanhu.io/g/sqlx"
)

type psql struct {
	db *sqlx.DB
}

// NewPsql wraps a database connection into a store.
func NewPsql(db *sqlx.DB) Store {
	if db == nil {
		return NewMemStore()
	}
	return &psql{db: db}
}

// InitPsql initializes the object store database.
func InitPsql(b *sqlx.DB) error {
	q := `create table objects (
		k varchar(255) primary key not null,
		v bytea not null
	)`
	_, err := b.X(q)
	return err
}

// DestroyPsql destroys the object store database.
func DestroyPsql(db *sqlx.DB) error {
	return sqlx.DestroyTable(db, "objects")
}

func (db *psql) Put(bs []byte) (string, error) {
	hash := hashutil.Hash(bs)

	q := `insert into objects (k, v)
		values ($1, $2) on conflict do nothing`
	if _, err := db.db.X(q, hash, bs); err != nil {
		return "", err
	}
	return hash, nil
}

func (db *psql) Get(hash string) ([]byte, error) {
	q := `select v from objects where k=$1`
	row := db.db.Q1(q, hash)

	var ret []byte
	if has, err := row.Scan(&ret); err != nil {
		return nil, err
	} else if !has {
		return nil, pathutil.NotExist(hash)
	}
	return ret, nil
}

func (db *psql) Has(hash string) (bool, error) {
	q := `select 1 from objects where k=$1`
	row := db.db.Q1(q, hash)
	var v int
	return row.Scan(&v)
}
