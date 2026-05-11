package timeutil

import (
	"time"

	"testing"
)

func TestNowFunc(t *testing.T) {
	now := time.Now()
	f := func() time.Time { return now }
	ts := NowFunc(f)()
	if !ts.Equal(now) {
		t.Errorf("read time got %q, want %q", ts, now)
	}

	later := NowFunc(nil)()
	if later.Before(now) {
		t.Errorf("time is reversing")
	}
}
