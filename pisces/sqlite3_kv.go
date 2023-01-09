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

type sqlite3KV struct {
	db    *sqlx.DB
	table string
}

func newSqlite3KV(db *sqlx.DB, table string) *sqlite3KV {
	return &sqlite3KV{db: db, table: table}
}

func (b *sqlite3KV) create() error {
	return Sqlite3CreateKV(b.db, b.table)
}

func (b *sqlite3KV) createMissing() error {
	return Sqlite3CreateKVMissing(b.db, b.table)
}

func (b *sqlite3KV) destroy() error {
	return Sqlite3Drop(b.db, b.table)
}

func (b *sqlite3KV) clear() error {
	q := fmt.Sprintf(`delete from %s`, b.table)
	_, err := b.db.X(q)
	return err
}

func (b *sqlite3KV) add(k, cls string, bs []byte) error {
	q := fmt.Sprintf(`insert into %s (k, c, v) values (?, ?, ?)`, b.table)
	_, err := b.db.X(q, k, cls, bs)
	return err
}

func (b *sqlite3KV) get(k string) ([]byte, error) {
	q := fmt.Sprintf(`select v from %s where k=?`, b.table)
	row := b.db.Q1(q, k)
	var bs []byte
	if has, err := row.Scan(&bs); err != nil {
		return nil, err
	} else if !has {
		return nil, notFound
	}
	return bs, nil
}

func (b *sqlite3KV) has(k string) (bool, error) {
	q := fmt.Sprintf(`select 1 from %s where k=?`, b.table)
	row := b.db.Q1(q, k)
	var i int
	if has, err := row.Scan(&i); err != nil {
		return false, err
	} else if !has {
		return false, nil
	}
	return true, nil
}

func (b *sqlite3KV) set(k string, bs []byte) error {
	q := fmt.Sprintf(`update %s set v=? where k=?`, b.table)
	res, err := b.db.X(q, bs, k)
	if err != nil {
		return err
	}
	return sqlResError(res)
}

func (b *sqlite3KV) setClass(k, cls string) error {
	q := fmt.Sprintf(`update %s set c=? where k=?`, b.table)
	res, err := b.db.X(q, cls, k)
	if err != nil {
		return err
	}
	return sqlResError(res)
}

func (b *sqlite3KV) remove(k string) error {
	q := fmt.Sprintf(`delete from %s where k=?`, b.table)
	res, err := b.db.X(q, k)
	if err != nil {
		return err
	}
	return sqlResError(res)
}

func (b *sqlite3KV) emplace(k, cls string, bs []byte) error {
	q := fmt.Sprintf(
		"insert into %s (k, v, c) values (?, ?, ?) "+
			"on conflict (k) do nothing", b.table,
	)
	_, err := b.db.X(q, k, bs, cls)
	return err
}

func (b *sqlite3KV) replace(k, cls string, bs []byte) error {
	q := fmt.Sprintf(
		"insert into %s (k, v, c) values (?, ?, ?) "+
			"on conflict (k) do update set v=excluded.v",
		b.table,
	)
	_, err := b.db.X(q, k, bs, cls)
	return err
}

func (b *sqlite3KV) appendBytes(k string, bs []byte) error {
	q := fmt.Sprintf(
		"insert into %s (k, v, c) values (?, ?, ?) "+
			"on conflict (k) do update set v = %s.v || excluded.v",
		b.table, b.table,
	)
	_, err := b.db.X(q, k, bs, "")
	return err
}

func (b *sqlite3KV) mutate(k string, f func(bs []byte) ([]byte, error)) error {
	tx, err := b.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var bs []byte
	q := fmt.Sprintf(`select v from %s where k=?`, b.table)
	row := tx.Q1(q, k)
	if has, err := row.Scan(&bs); err != nil {
		return err
	} else if !has {
		return notFound
	}

	newBytes, err := f(bs)
	if err != nil {
		return err
	}

	q = fmt.Sprintf(`update %s set v=? where k=?`, b.table)
	res, err := tx.X(q, newBytes, k)
	if err != nil {
		return err
	}
	if err := sqlResError(res); err != nil {
		return err
	}
	return tx.Commit()
}

func (b *sqlite3KV) count() (int64, error) {
	row := b.db.Q1(fmt.Sprintf(`select count(1) from %s`, b.table))
	var v int64
	if has, err := row.Scan(&v); err != nil {
		return 0, err
	} else if !has {
		return 0, errcode.Internalf("count returns nothing")
	}
	return v, nil
}

func (b *sqlite3KV) walk(f WalkFunc) error {
	q := fmt.Sprintf(`select k, c, v from %s order by k`, b.table)
	rows, err := b.db.Q(q)
	if err != nil {
		return err
	}
	defer rows.Close()
	return sqlIterRows(rows, f)
}

func (b *sqlite3KV) walkClass(cls string, f WalkFunc) error {
	q := fmt.Sprintf(`select k, c, v from %s where c=? order by k`, b.table)
	rows, err := b.db.Q(q, cls)
	if err != nil {
		return err
	}
	defer rows.Close()
	return sqlIterRows(rows, f)
}

func (b *sqlite3KV) walkPartial(p *KVPartial, f WalkFunc) error {
	q := fmt.Sprintf(
		"select k, c, v from %s order by k %s limit %d offset %d",
		b.table, sqlOrderStr(p.Desc), p.N, p.Offset,
	)
	rows, err := b.db.Q(q)
	if err != nil {
		return err
	}
	defer rows.Close()
	return sqlIterRows(rows, f)
}

func (b *sqlite3KV) walkPartialClass(
	cls string, p *KVPartial, f WalkFunc,
) error {
	q := fmt.Sprintf(
		"select k, c, v from %s where c=? "+
			"order by k %s limit %d offset %d",
		b.table, sqlOrderStr(p.Desc), p.N, p.Offset,
	)
	rows, err := b.db.Q(q, cls)
	if err != nil {
		return err
	}
	defer rows.Close()
	return sqlIterRows(rows, f)
}

func (b *sqlite3KV) ops() *KVOps {
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

// NewSqlite3KV creates a new sqlite3 backed unordered key-value pair storage.
func NewSqlite3KV(db *sqlx.DB, table string) *KV {
	return newKV(newSqlite3KV(db, table).ops())
}

// NewOrderedSqlite3KV creates a new sqlite3 backed ordered key-value pair
// storage.
func NewOrderedSqlite3KV(db *sqlx.DB, table string) *KV {
	return newOrderedKV(newSqlite3KV(db, table).ops())
}
