package timeutil

import (
	"testing"

	"time"
)

func TestTimestamp(t *testing.T) {
	now := time.Now()
	nanos := now.UnixNano()
	nanos2 := NewTimestamp(now).Time().UnixNano()
	if nanos2 != nanos {
		t.Errorf(
			"timestamp roundtrip failed: %s: %d != %d",
			now, nanos, nanos2,
		)
	}

	got := Time(nil)
	var zeroTime time.Time
	if !got.Equal(zeroTime) {
		t.Errorf("nil timestamp got %q", got)
	}
}
