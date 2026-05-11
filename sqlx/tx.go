package sqlx

import (
	"database/sql"
)

// Tx wraps a transaction
type Tx struct {
	*sql.Tx
	*wrap
}
