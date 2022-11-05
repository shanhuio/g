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
	"testing"

	"time"
)

func TestGapper(t *testing.T) {
	now := time.Now()
	ticker := newGapperNow(time.Second, now)

	type testPoint struct {
		d    time.Duration
		want bool
	}
	for i, test := range []*testPoint{
		{d: time.Duration(0), want: false},
		{d: time.Second, want: true},
		{d: time.Second, want: false},
		{d: time.Second + time.Second/2, want: false},
		{d: 2*time.Second + time.Second/2, want: true},
		{d: 2*time.Second + time.Second*7/10, want: false},
		{d: 3 * time.Second, want: false},
		{d: 3*time.Second + time.Second/2, want: true},
	} {
		got := ticker.check(now.Add(test.d))
		if got != test.want {
			t.Errorf(
				"check #%d with %s, got %t, want %t",
				i, test.d, got, test.want,
			)
		}
	}
}
