package sqlx

import (
	"database/sql"
)

// Row is a result with the row and the query.
type Row struct {
	Query string
	*sql.Row
}

// Scan scans a row into values.
func (r *Row) Scan(dest ...interface{}) (bool, error) {
	err := r.Row.Scan(dest...)
	if err == nil {
		return true, nil
	}
	if err == sql.ErrNoRows {
		return false, nil
	}
	return false, Error(r.Query, err)
}

// Rows is a result with the rows and the query.
type Rows struct {
	Query string
	*sql.Rows
}

// Close closes the rows result.
func (r *Rows) Close() error {
	return Error(r.Query, r.Rows.Close())
}

// Err returns the error.
func (r *Rows) Err() error {
	return Error(r.Query, r.Rows.Err())
}

// Scan scans a row into values.
func (r *Rows) Scan(dest ...interface{}) error {
	return Error(r.Query, r.Rows.Scan(dest...))
}
