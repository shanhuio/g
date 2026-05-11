package idutil

import (
	"testing"
)

func TestShortId(t *testing.T) {
	for _, test := range []struct {
		id   string
		want string
	}{
		{"", ""},
		{"123", "123"},
		{"1234567", "1234567"},
		{"12345678", "1234567"},
		{"1234567890", "1234567"},
		{"\\\\", ""},
		{"汉字？？", ""},
	} {
		got := Short(test.id)
		if got != test.want {
			t.Errorf(
				"Short id string for %q: got %q, want %q",
				test.id, got, test.want,
			)
		}
	}
}
