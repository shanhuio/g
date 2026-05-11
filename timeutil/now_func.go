package timeutil

import (
	"time"
)

// NowFunc returns f if f is not nil, otherwise returns time.Now
func NowFunc(f func() time.Time) func() time.Time {
	if f != nil {
		return f
	}
	return time.Now
}

// ReadTime runs the function and returns the time if f is not null, or returns
// time.Now() if f is null.
func ReadTime(f func() time.Time) time.Time {
	if f == nil {
		return time.Now()
	}
	return f()
}

// ReadTimestamp runs the function and returns the time as a Timestamp,
// or return the result of time.Now() as a timestamp.
func ReadTimestamp(f func() time.Time) *Timestamp {
	return NewTimestamp(ReadTime(f))
}
