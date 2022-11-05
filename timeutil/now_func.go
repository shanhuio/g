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
