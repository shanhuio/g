package sqlx

import (
	"database/sql"
)

type conn interface {
	Exec(q string, args ...interface{}) (sql.Result, error)
	Query(q string, args ...interface{}) (*sql.Rows, error)
	QueryRow(q string, args ...interface{}) *sql.Row
}
