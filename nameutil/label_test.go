package nameutil

import (
	"testing"
)

func TestCheckLabel(t *testing.T) {
	for _, name := range []string{
		"normal-label",
		"normal-long-label",
		"s",   // single letter
		"888", // number only
	} {
		if err := CheckLabel(name); err != nil {
			t.Errorf("%q is a valid label, but got error %q", name, err)
		}
	}

	for _, name := range []string{
		"",
		"A-very-very-very-very-very-very-very-very–very-very-long",
		"-label",
		"label-",
		"label--label",
		"label&label",
	} {
		if CheckLabel(name) == nil {
			t.Errorf("%q is an invalid label, but got no checking error", name)
		}
	}
}
