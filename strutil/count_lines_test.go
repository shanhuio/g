package strutil

import (
	"testing"
)

func TestCountLines(t *testing.T) {
	for _, test := range []struct {
		bs   []byte
		want int
	}{
		{nil, 0},
		{make([]byte, 0), 0},
		{[]byte("abcd"), 1},
		{[]byte("\n"), 1},
		{[]byte(" \n"), 1},
		{[]byte("\n\n"), 2},
		{[]byte("\n\nabc"), 3},
	} {
		got := CountLines(test.bs)
		if test.want != got {
			t.Errorf(
				"CountLines(%q): got %d, want %d",
				string(test.bs), got, test.want,
			)
		}
	}
}
