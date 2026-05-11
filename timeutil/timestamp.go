package timeutil

import (
	"time"
)

// Timestamp is a struct to record a UTC timestamp.
// It is designed to be directly usable in Javascript.
type Timestamp struct {
	Sec  int64
	Nano int64 `json:",omitempty"`
}

// Time returns the time of this timestamp in UTC.
func (t *Timestamp) Time() time.Time {
	return time.Unix(t.Sec, t.Nano).UTC()
}

// Clone clones the timestamp.
func (t *Timestamp) Clone() *Timestamp {
	cp := *t
	return &cp
}

// NewTimestamp creates a new timestamp from the given time.
func NewTimestamp(t time.Time) *Timestamp {
	sec, nano := secNano(t.UnixNano())
	return &Timestamp{
		Sec:  sec,
		Nano: nano,
	}
}

// TimestampNow creates a time stamp of the time now.
func TimestampNow() *Timestamp {
	return NewTimestamp(time.Now())
}

// Time converts timestamp to time.Time .
func Time(ts *Timestamp) time.Time {
	if ts == nil {
		var zero time.Time
		return zero
	}
	return ts.Time()
}

// CopyTimestamp copies a timestamp.
func CopyTimestamp(ts *Timestamp) *Timestamp {
	if ts == nil {
		return nil
	}
	return ts.Clone()
}
