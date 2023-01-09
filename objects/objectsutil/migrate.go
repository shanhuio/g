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

package objectsutil

import (
	"bytes"
	"fmt"

	"shanhu.io/pub/objects"
	"shanhu.io/pub/sqlx"
)

// DumpPsqlDB dumps a psql database into a streaming object store.
func DumpPsqlDB(db *sqlx.DB, objs objects.Objects) error {
	q := `select k, v from objects`
	var k string
	var v []byte
	rows, err := db.Q(q)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&k, &v); err != nil {
			return err
		}

		buf := bytes.NewReader(v)
		newKey, err := objs.Create(buf)
		if err != nil {
			return err
		}
		if newKey != k {
			return fmt.Errorf("key changed from %q to %q", k, newKey)
		}
	}

	return rows.Close()
}
