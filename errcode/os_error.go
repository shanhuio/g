package errcode

import (
	stderrcode "shanhu.io/std/errcode"
)

// FromOS converts os package errors into errcode errors.
func FromOS(err error) error { return stderrcode.FromOS(err) }
