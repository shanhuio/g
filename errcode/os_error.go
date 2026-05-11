package errcode

import (
	"os"
)

// FromOS converts os package errors into errcode errors.
func FromOS(err error) error {
	if err == nil {
		return err
	}
	if os.IsNotExist(err) {
		return Add(NotFound, err)
	}
	if os.IsPermission(err) {
		return Add(Unauthorized, err)
	}
	if os.IsTimeout(err) {
		return Add(TimeOut, err)
	}
	return err
}
