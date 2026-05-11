package pisces

import (
	"testing"

	"os"
	"path/filepath"

	_ "modernc.org/sqlite" // sqlite db driver
	"shanhu.io/g/sqlx"
)

func TestSqlite3KV(t *testing.T) {
	dir, err := os.MkdirTemp("", "pisces")
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
