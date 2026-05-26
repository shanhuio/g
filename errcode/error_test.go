package errcode

import (
	"testing"

	stderrcode "shanhu.io/std/errcode"
)

func TestCommonError(t *testing.T) {
	for _, test := range []struct {
		f   func(err error) bool
		err error
	}{
		{IsNotFound, NotFoundf("not found")},
		{IsInvalidArg, InvalidArgf("invalid arg")},
		{IsUnauthorized, Unauthorizedf("unauthorized")},
		{IsInternal, Internalf("internal")},
		{IsTimeOut, TimeOutf("time out")},
	} {
		if !test.f(test.err) {
			t.Errorf("test failed for error: %s", test.err)
		}
	}
}

func TestStdNotFoundIsNotFound(t *testing.T) {
	err := stderrcode.NotFoundf("missing")
	if !IsNotFound(err) {
		t.Errorf("IsNotFound(%v) = false, want true", err)
	}
	if got := Of(err); got != NotFound {
		t.Errorf("Of(%v) = %q, want %q", err, got, NotFound)
	}
}

func TestNotFoundIsStdNotFound(t *testing.T) {
	err := NotFoundf("missing")
	if !stderrcode.IsNotFound(err) {
		t.Errorf("stderrcode.IsNotFound(%v) = false, want true", err)
	}
	if got := stderrcode.Of(err); got != stderrcode.NotFound {
		t.Errorf("stderrcode.Of(%v) = %q, want %q", err, got, stderrcode.NotFound)
	}
}
