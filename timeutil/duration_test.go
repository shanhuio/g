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
	"testing"

	"time"
)

func TestDuration(t *testing.T) {
	for _, d := range []time.Duration{
		time.Duration(0),
		time.Nanosecond,
		time.Second,
		time.Minute,
		time.Hour,
		time.Duration(123456789),
	} {
		d2 := TimeDuration(NewDuration(d))
		if d2 != d {
			t.Errorf("cycling duration %s got %s", d, d2)
		}
	}

	got := TimeDuration(nil)
	if got != time.Duration(0) {
		t.Errorf("want 0 duration, got %s", got)
	}
}
