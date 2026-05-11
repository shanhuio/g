package timeutil

import (
	"testing"

	"time"
)

func TestDuration(t *testing.T) {
	for _, d := range []time.Duration{
		time.Duration(0),
		time.Nanosecond,
		time.Second,
		time.Minute,
		time.Hour,
		time.Duration(123456789),
	} {
		d2 := TimeDuration(NewDuration(d))
		if d2 != d {
			t.Errorf("cycling duration %s got %s", d, d2)
		}
	}

	got := TimeDuration(nil)
	if got != time.Duration(0) {
		t.Errorf("want 0 duration, got %s", got)
	}
}
