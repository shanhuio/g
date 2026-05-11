package certutil

import (
	"testing"

	"time"
)

func TestGapper(t *testing.T) {
	now := time.Now()
	ticker := newGapperNow(time.Second, now)

	type testPoint struct {
		d    time.Duration
		want bool
	}
	for i, test := range []*testPoint{
		{d: time.Duration(0), want: false},
		{d: time.Second, want: true},
		{d: time.Second, want: false},
		{d: time.Second + time.Second/2, want: false},
		{d: 2*time.Second + time.Second/2, want: true},
		{d: 2*time.Second + time.Second*7/10, want: false},
		{d: 3 * time.Second, want: false},
		{d: 3*time.Second + time.Second/2, want: true},
	} {
		got := ticker.check(now.Add(test.d))
		if got != test.want {
			t.Errorf(
				"check #%d with %s, got %t, want %t",
				i, test.d, got, test.want,
			)
		}
	}
}
