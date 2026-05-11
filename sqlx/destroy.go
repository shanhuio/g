package sqlx

// DestroyTable is a helper for destroying a table.
func DestroyTable(db *DB, t string) error {
	_, err := db.X("drop table if exists " + t)
	return err
}

// DestroyTables is a helper for destroying a set of tables.
func DestroyTables(db *DB, ts ...string) error {
	for _, t := range ts {
		if err := DestroyTable(db, t); err != nil {
			return err
		}
	}
	return nil
}
