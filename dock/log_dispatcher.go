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

package dock

import (
	"encoding/binary"
	"fmt"
	"io"
)

type logDispatcher struct {
	Stdout io.Writer
	Stderr io.Writer
}

func (d *logDispatcher) pipe(r io.Reader) error {
	header := make([]byte, 8)
	for {
		if _, err := io.ReadFull(r, header); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		var out io.Writer
		stream := header[0]
		switch stream {
		case 0, 1: // stdout
			out = d.Stdout
		case 2:
			out = d.Stderr
		default:
			return fmt.Errorf("invalid stream %d in header", stream)
		}

		n := binary.BigEndian.Uint32(header[4:8])
		if _, err := io.CopyN(out, r, int64(n)); err != nil {
			return err
		}
	}
}
