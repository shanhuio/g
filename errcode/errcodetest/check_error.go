package errcodetest

import (
	"testing"

	"shanhu.io/g/errcode"
)

// CheckError checks whether the error has expected error code for tests
func CheckError(t *testing.T, err error, code, message string) {
	t.Helper()
	if err == nil {
		t.Errorf("got nil, want err: %s", message)
	}
	if errcode.Of(err) != code {
		t.Errorf(
			"got error code %s, want %s: %s",
			errcode.Of(err),
			code,
			message,
		)
	}
}
