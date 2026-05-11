package signer

import (
	"time"
)

const timestampLen = 8

func now(f func() time.Time) time.Time {
	if f == nil {
		return time.Now()
	}
	return f()
}

func inWindow(t, tnow time.Time, w time.Duration) bool {
	tstart := tnow.Add(-w)
	tend := tnow.Add(w)
	return t.After(tstart) && t.Before(tend)
}
