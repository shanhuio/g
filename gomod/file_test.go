package gomod

import (
	"testing"
)

func TestModulePath(t *testing.T) {
	for _, test := range []struct {
		content, mod string
	}{
		{`module shanhu.io/pub`, "shanhu.io/pub"},
		{"  module    shanhu.io/pub\t\t\t\n\nextra", "shanhu.io/pub"},
		{`module "shanhu.io/g/v1"`, "shanhu.io/g/v1"},
		{`module "shanhu.io/pub"`, "shanhu.io/pub"},
		{"// comment\nmodule x // tail\nnext line", "x"},
		{"module `x` // tail", "x"},
	} {
		got, err := modulePath([]byte(test.content))
		if err != nil {
			t.Errorf("modulePath(%q) got error: %s", test.content, err)
		} else if got != test.mod {
			t.Errorf(
				"modulePath(%q), want %q, got %q",
				test.content, test.mod, got,
			)
		}
	}
}
