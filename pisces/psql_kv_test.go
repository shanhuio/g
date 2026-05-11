package pisces

import (
	"testing"

	"os"

	_ "github.com/lib/pq"
	"shanhu.io/g/sqlx"
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
