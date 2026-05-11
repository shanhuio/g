package sqlx

import (
	"database/sql"
)

type wrap struct {
	conn
}

// X executes a query string.
func (w *wrap) X(q string, args ...interface{}) (sql.Result, error) {
	res, err := w.conn.Exec(q, args...)
	return res, Error(q, err)
}

// Q1 executes a query string that expects one single row as the return.
func (w *wrap) Q1(q string, args ...interface{}) *Row {
	return &Row{
		Query: q,
		Row:   w.conn.QueryRow(q, args...),
	}
}

// Q queries the database with a query string.
func (w *wrap) Q(q string, args ...interface{}) (*sql.Rows, error) {
	res, err := w.conn.Query(q, args...)
	if err != nil {
		return nil, Error(q, err)
	}
	return res, nil
}
