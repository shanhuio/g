package pathutil

import (
	"testing"
)

func TestIsParent(t *testing.T) {
	for _, test := range []struct {
		short, long string
		want        bool
	}{
		{"a", "b", false},
		{"a", "a", true},
		{"a/b", "a", false},
		{"a", "a/b", true},
		{"a", "ab", false},
		{"a", "ab/c", false},
		{"a/b", "a/b/c", true},
	} {
		got := IsParent(test.short, test.long)
		if got != test.want {
			t.Errorf(
				"IsParent(%q, %q), want %v, got %v",
				test.short, test.long, test.want, got,
			)
		}
	}
}
