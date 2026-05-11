package goload

import (
	"fmt"
)

type scanError struct {
	dir string
	err error
}

func (err *scanError) Error() string {
	return fmt.Sprintf("scan %q: %s", err.dir, err.err)
}
