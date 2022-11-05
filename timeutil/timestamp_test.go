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

func TestTimestamp(t *testing.T) {
	now := time.Now()
	nanos := now.UnixNano()
	nanos2 := NewTimestamp(now).Time().UnixNano()
	if nanos2 != nanos {
		t.Errorf(
			"timestamp roundtrip failed: %s: %d != %d",
			now, nanos, nanos2,
		)
	}

	got := Time(nil)
	var zeroTime time.Time
	if !got.Equal(zeroTime) {
		t.Errorf("nil timestamp got %q", got)
	}
}
