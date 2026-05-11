package testutil

import (
	"log"
)

// Failer is an interface to mark that a test has failed.
type Failer interface {
	Fail()
}

// PanicFailer fails with a log.Fatal call.
type PanicFailer struct{}

// Fail fails with a log.Fatal call.
func (pf *PanicFailer) Fail() {
	log.Fatal("test failed")
}
