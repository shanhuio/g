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
	"testing"

	"io/ioutil"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite" // sqlite db driver
	"shanhu.io/pub/sqlx"
)

func TestSqlite3KV(t *testing.T) {
	dir, err := ioutil.TempDir("", "pisces")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	dbFile := filepath.Join(dir, "db")
	db, err := sqlx.OpenSqlite3(dbFile)
	if err != nil {
		t.Fatal(err)
	}

	const testTable = "testkv"

	for _, test := range kvTestSuite {
		t.Log(test.name)

		if err := Sqlite3DropExist(db, testTable); err != nil {
			t.Fatal(err)
		}
		if err := Sqlite3CreateKV(db, testTable); err != nil {
			t.Fatal(err)
		}

		var kv *KV
		if !test.ordered {
			kv = NewSqlite3KV(db, testTable)
		} else {
			kv = NewOrderedSqlite3KV(db, testTable)
		}
		test.f(t, kv)
	}
}
