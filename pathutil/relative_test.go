package pathutil

import (
	"testing"
)

func TestRelative(t *testing.T) {
	for _, test := range []struct {
		base, full string
		want       string
	}{
		{"a", "b", ""},
		{"a", "a/b", "b"},
		{"a", "a", "."},
		{"a", "a/b/c", "b/c"},
		{"a", "ab/c", ""},
	} {
		got := Relative(test.base, test.full)
		if got != test.want {
			t.Errorf(
				"Relative(%q, %q), want %q, got %q",
				test.base, test.full, test.want, got,
			)
		}
	}
}

func TestDotRelative(t *testing.T) {
	for _, test := range []struct {
		base, full string
		want       string
	}{
		{"a", "b", "b"},
		{"a", "a/b", "./b"},
		{"a", "a", "."},
		{"a", "a/b/c", "./b/c"},
		{"a", "ab/c", "ab/c"},
	} {
		got := DotRelative(test.base, test.full)
		if got != test.want {
			t.Errorf(
				"DotRelative(%q, %q), want %q, got %q",
				test.base, test.full, test.want, got,
			)
		}
	}
}
