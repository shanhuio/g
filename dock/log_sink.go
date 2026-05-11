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
