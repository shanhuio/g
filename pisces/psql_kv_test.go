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

	"os"

	_ "github.com/lib/pq"
	"shanhu.io/pub/sqlx"
)

func TestPsqlKV(t *testing.T) {
	dbPath := os.Getenv("GOTESTPSQL")
	if dbPath == "" {
		t.Log("GOTESTPSQL not set, skipping.")
		return
	}
	db, err := sqlx.OpenPsql(dbPath)
	if err != nil {
		t.Fatal(err)
	}

	const testTable = "testkv"

	for _, test := range kvTestSuite {
		t.Log(test.name)

		if err := PsqlDropExist(db, testTable); err != nil {
			t.Fatal(err)
		}

		if err := PsqlCreateKV(db, testTable); err != nil {
			t.Fatal(err)
		}

		var kv *KV
		if !test.ordered {
			kv = NewPsqlKV(db, testTable)
		} else {
			kv = NewOrderedPsqlKV(db, testTable)
		}
		test.f(t, kv)
	}

	if err := PsqlDropExist(db, testTable); err != nil {
		t.Fatal(err)
	}
}
