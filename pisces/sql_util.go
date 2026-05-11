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
