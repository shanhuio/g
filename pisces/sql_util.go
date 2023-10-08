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
	"database/sql"

	"shanhu.io/g/sqlx"
)

func sqlResError(res sql.Result) error {
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return notFound
	}
	if n != 1 {
		return multiAffected
	}
	return nil
}

func sqlIterRows(rows *sql.Rows, f WalkFunc) error {
	for rows.Next() {
		var k, cls string
		var bs []byte
		if err := rows.Scan(&k, &cls, &bs); err != nil {
			return err
		}
		if err := f(k, cls, bs); err != nil {
			return err
		}
	}
	return rows.Close()
}

func sqlOrderStr(desc bool) string {
	if desc {
		return "desc"
	}
	return "asc"
}

func tableDriver(db *sqlx.DB) string {
	if db == nil {
		return ""
	}
	return db.Driver()
}
