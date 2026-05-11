package strutil

import (
	"testing"
)

func TestDefault(t *testing.T) {
	want := "some string"
	for _, test := range []struct {
		input, def, want string
	}{
		{"", want, want},
		{want, "default", want},
	} {
		got := Default(test.input, test.def)
		if got != want {
			t.Errorf(
				"Default(%q, %q), want %q, got %q",
				test.input, test.def, test.want, got,
			)
		}
	}
}
