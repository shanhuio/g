package pisces

import (
	"fmt"

	"shanhu.io/g/sqlx"
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
