package sqlx

// RunTasks runs a series of tasks, and returns error on the first one that
// fails.
func RunTasks(db *DB, funcs ...func(db *DB) error) error {
	for _, f := range funcs {
		if err := f(db); err != nil {
			return err
		}
	}
	return nil
}
