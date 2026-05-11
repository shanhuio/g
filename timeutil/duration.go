package timeutil

import (
	"time"
)

func secNano(nano int64) (int64, int64) {
	sec := nano / 1e9
	nano -= sec * 1e9
	if nano < 0 {
		nano += 1e9
		sec--
	}
	return sec, nano
}

// Duration is a struct to record a duration.
// It is designed to be directly usable in Javascript.
type Duration struct {
	Sec  int64
	Nano int64 `json:",omitempty"`
}

// NewDuration creates a new duration of d.
func NewDuration(d time.Duration) *Duration {
	sec, nano := secNano(d.Nanoseconds())
	return &Duration{
		Sec:  sec,
		Nano: nano,
	}
}

// Duration converts to time.Duration type.
func (d *Duration) Duration() time.Duration {
	sec := time.Duration(d.Sec) * time.Second
	nano := time.Duration(d.Nano) * time.Nanosecond
	return sec + nano
}

// TimeDuration converts d to time.Duration. If d is nil, 0 duration is
// returned.
func TimeDuration(d *Duration) time.Duration {
	if d == nil {
		return time.Duration(0)
	}
	return d.Duration()
}
