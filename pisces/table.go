package pisces

// Table contains common operations for creating and destroying a table.
type Table interface {
	// Create creates table.
	Create() error

	// CreateMissing creates the table if it is missing.
	CreateMissing() error

	// Destroy destroys the table.
	Destroy() error
}
