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

package bufpool

import (
	"testing"
)

func TestBytes(t *testing.T) {
	const size = 1024
	p := NewBytes(1024)
	buf := p.Get()
	if len(buf) != size {
		t.Errorf("buf got size %d, want %d", len(buf), size)
	}

	p.Put(buf)
}

func TestBytes_defaultSize(t *testing.T) {
	p := NewBytes(0)
	buf := p.Get()
	if len(buf) != DefaultBytesSize {
		t.Errorf("buf got size %d, want %d", len(buf), DefaultBytesSize)
	}
}
