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

	"shanhu.io/g/errcode"
	"shanhu.io/g/sqlx"
)

type psqlKV struct {
	db    *sqlx.DB
	table string
}

func newPsqlKV(db *sqlx.DB, table string) *psqlKV {
	return &psqlKV{db: db, table: table}
}

func (b *psqlKV) create() error {
	return PsqlCreateKV(b.db, b.table)
}

func (b *psqlKV) createMissing() error {
	return PsqlCreateKVMissing(b.db, b.table)
}

func (b *psqlKV) destroy() error {
	return PsqlDrop(b.db, b.table)
}

func (b *psqlKV) clear() error {
	q := fmt.Sprintf(`truncate table %s`, b.table)
	_, err := b.db.X(q)
	return err
}

func (b *psqlKV) add(k, cls string, bs []byte) error {
	q := fmt.Sprintf(`insert into %s (k, c, v) values ($1, $2, $3)`, b.table)
	_, err := b.db.X(q, k, cls, bs)
	return err
}

func (b *psqlKV) get(k string) ([]byte, error) {
	q := fmt.Sprintf(`select v from %s where k=$1`, b.table)
	row := b.db.Q1(q, k)
	var bs []byte
	if has, err := row.Scan(&bs); err != nil {
		return nil, err
	} else if !has {
		return nil, notFound
	}
	return bs, nil
}

func (b *psqlKV) has(k string) (bool, error) {
	q := fmt.Sprintf(`select 1 from %s where k=$1`, b.table)
	row := b.db.Q1(q, k)
	var i int
	if has, err := row.Scan(&i); err != nil {
		return false, err
	} else if !has {
		return false, nil
	}
	return true, nil
}

func (b *psqlKV) set(k string, bs []byte) error {
	q := fmt.Sprintf(`update %s set v=$1 where k=$2`, b.table)
	res, err := b.db.X(q, bs, k)
	if err != nil {
		return err
	}
	return sqlResError(res)
}

func (b *psqlKV) setClass(k, cls string) error {
	q := fmt.Sprintf(`update %s set c=$1 where k=$2`, b.table)
	res, err := b.db.X(q, cls, k)
	if err != nil {
		return err
	}
	return sqlResError(res)
}

func (b *psqlKV) remove(k string) error {
	q := fmt.Sprintf(`delete from %s where k=$1`, b.table)
	res, err := b.db.X(q, k)
	if err != nil {
		return err
	}
	return sqlResError(res)
}

func (b *psqlKV) emplace(k, cls string, bs []byte) error {
	q := fmt.Sprintf(
		"insert into %s (k, v, c) values ($1, $2, $3) "+
			"on conflict (k) do nothing", b.table,
	)
	_, err := b.db.X(q, k, bs, cls)
	return err
}

func (b *psqlKV) replace(k, cls string, bs []byte) error {
	q := fmt.Sprintf(
		"insert into %s (k, v, c) values ($1, $2, $3) "+
			"on conflict (k) do update set v=excluded.v",
		b.table,
	)
	_, err := b.db.X(q, k, bs, cls)
	return err
}

func (b *psqlKV) appendBytes(k string, bs []byte) error {
	q := fmt.Sprintf(
		"insert into %s (k, v, c) values ($1, $2, $3) "+
			"on conflict (k) do update set v = %s.v || excluded.v",
		b.table, b.table,
	)
	_, err := b.db.X(q, k, bs, "")
	return err
}

func (b *psqlKV) mutate(k string, f func(bs []byte) ([]byte, error)) error {
	tx, err := b.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var bs []byte
	q := fmt.Sprintf(`select v from %s where k=$1`, b.table)
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

	q = fmt.Sprintf(`update %s set v=$1 where k=$2`, b.table)
	res, err := tx.X(q, newBytes, k)
	if err != nil {
		return err
	}
	if err := sqlResError(res); err != nil {
		return err
	}
	return tx.Commit()
}

func (b *psqlKV) count() (int64, error) {
	row := b.db.Q1(fmt.Sprintf(`select count(1) from %s`, b.table))
	var v int64
	if has, err := row.Scan(&v); err != nil {
		return 0, err
	} else if !has {
		return 0, errcode.Internalf("count returns nothing")
	}
	return v, nil
}

func (b *psqlKV) walk(f WalkFunc) error {
	q := fmt.Sprintf(`select k, c, v from %s order by k`, b.table)
	rows, err := b.db.Q(q)
	if err != nil {
		return err
	}
	defer rows.Close()
	return sqlIterRows(rows, f)
}

func (b *psqlKV) walkClass(cls string, f WalkFunc) error {
	q := fmt.Sprintf(`select k, c, v from %s where c=$1 order by k`, b.table)
	rows, err := b.db.Q(q, cls)
	if err != nil {
		return err
	}
	defer rows.Close()
	return sqlIterRows(rows, f)
}

func (b *psqlKV) walkPartial(p *KVPartial, f WalkFunc) error {
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

func (b *psqlKV) walkPartialClass(
	cls string, p *KVPartial, f WalkFunc,
) error {
	q := fmt.Sprintf(
		"select k, c, v from %s where c=$1 "+
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

func (b *psqlKV) ops() *KVOps {
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

// NewPsqlKV creates a new postgresql backed unordered key-value pair storage.
func NewPsqlKV(db *sqlx.DB, table string) *KV {
	return newKV(newPsqlKV(db, table).ops())
}

// NewOrderedPsqlKV creates a new postgresql backed ordered key-value pair
// storage.
func NewOrderedPsqlKV(db *sqlx.DB, table string) *KV {
	return newOrderedKV(newPsqlKV(db, table).ops())
}
