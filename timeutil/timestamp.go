// Copyright (C) 2022  Shanhu Tech Inc.
//
// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU Affero General Public License as published by the
// Free Software Foundation, either version 3 of the License, or (at your
// option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.  See the GNU Affero General Public License
// for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
