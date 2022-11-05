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

package counting

import (
	"sync/atomic"
)

// Counter is an atomic counter
type Counter struct {
	count int64
}

// NewCounter creates a Counter instance and initializes it to 0.
func NewCounter() *Counter {
	return &Counter{0}
}

// Add adds an int to the counter
func (c *Counter) Add(i int64) {
	atomic.AddInt64(&c.count, int64(i))
}

// Count returns count
func (c *Counter) Count() int64 {
	return atomic.LoadInt64(&c.count)
}

// Value is a alias of Count
func (c *Counter) Value() int64 {
	return c.Count()
}

// Set sets the count to the new value and returns previous value.
func (c *Counter) Set(n int64) int64 {
	return atomic.SwapInt64(&c.count, n)
}

// Reset change count to 0, returns previous count
func (c *Counter) Reset() int64 {
	return c.Set(0)
}
