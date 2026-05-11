package errcode

import (
	"testing"
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
