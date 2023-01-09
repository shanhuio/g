// Copyright (C) 2023  Shanhu Tech Inc.
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

	"crypto/rand"
	"time"
)

func TestChallenge(t *testing.T) {
	now := time.Now()

	r := rand.Reader
	ch, err := NewChallenge(now, r)
	if err != nil {
		t.Fatal(err)
	}

	got := Time(ch.T)
	if !now.Equal(got) {
		t.Errorf("got timestamp %q, want %q", got, now)
	}
	if ch.N == "" {
		t.Errorf("nounce is empty")
	}

	ch2, err := NewChallenge(now, r)
	if err != nil {
		t.Fatal("get second challenge: ", err)
	}

	if ch2.N == ch.N {
		t.Errorf("got same nounce: %q", ch.N)
	}
}
