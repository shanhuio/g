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
