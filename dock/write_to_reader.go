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
	"io"
)

type writeToReader struct {
	wt   io.WriterTo
	r    *io.PipeReader
	w    *io.PipeWriter
	err  error
	done chan struct{}
}

func newWriteToReader(wt io.WriterTo) *writeToReader {
	r, w := io.Pipe()

	ret := &writeToReader{
		r:    r,
		w:    w,
		wt:   wt,
		done: make(chan struct{}),
	}

	go func() {
		defer close(ret.done)
		ret.err = ret.pipe()
	}()
	return ret
}

func (r *writeToReader) Close() error {
	return r.r.Close()
}

func (r *writeToReader) Join() error {
	<-r.done
	return r.err
}

func (r *writeToReader) pipe() error {
	if _, err := r.wt.WriteTo(r.w); err != nil {
		return err
	}
	return r.w.Close()
}

func (r *writeToReader) Read(buf []byte) (int, error) {
	return r.r.Read(buf)
}
