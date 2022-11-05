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

package certutil

import (
	"sync"
	"time"
)

// gapper provides a check function that returns true where the first true is
// after the first specified time point, and every two other adjacent true has
// a gap of no less than period.
type gapper struct {
	mu     sync.Mutex
	next   time.Time
	period time.Duration
}

func newGapper(period time.Duration, first time.Time) *gapper {
	return &gapper{
		next:   first,
		period: period,
	}
}

func newGapperNow(period time.Duration, now time.Time) *gapper {
	return newGapper(period, now.Add(period))
}

// check returns true if now is after the first specified time point if this is
// the first true, or if now is not before period after the last time check()
// returns true.
func (t *gapper) check(now time.Time) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !now.Before(t.next) {
		t.next = now.Add(t.period)
		return true
	}
	return false
}
