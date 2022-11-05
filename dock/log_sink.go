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
	"os"
)

type logSink struct {
	io.WriteCloser
	r    *io.PipeReader
	w    *io.PipeWriter
	d    *logDispatcher
	done chan error
}

func newLogSink(stdout, stderr io.Writer) *logSink {
	if stdout == nil {
		stdout = os.Stdout
	}
	if stderr == nil {
		stderr = os.Stderr
	}

	r, w := io.Pipe()
	d := &logDispatcher{
		Stdout: stdout,
		Stderr: stderr,
	}

	ret := &logSink{
		WriteCloser: w, r: r, w: w, d: d,
		done: make(chan error, 1),
	}

	go func() {
		ret.done <- ret.d.pipe(ret.r)
	}()
	return ret
}

func newStdLogSink() *logSink {
	return newLogSink(nil, nil)
}

func (s *logSink) waitDone() error {
	if err := s.w.Close(); err != nil {
		return err
	}
	return <-s.done
}
