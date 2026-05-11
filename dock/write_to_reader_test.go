package dock

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

type errWriterTo struct{ err error }

func (e *errWriterTo) WriteTo(io.Writer) (int64, error) {
	return 0, e.err
}

func TestWriteToReader(t *testing.T) {
	for _, test := range []struct {
		name string
		in   []byte
	}{
		{name: "nonempty", in: []byte("hello, world")},
		{name: "empty", in: nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			wr := newWriteToReader(bytes.NewReader(test.in))
			defer wr.Close()

			got, err := io.ReadAll(wr)
			if err != nil {
				t.Fatalf("ReadAll: %v", err)
			}
			if !bytes.Equal(got, test.in) {
				t.Errorf("got %q, want %q", got, test.in)
			}
			if err := wr.Join(); err != nil {
				t.Errorf("Join: %v", err)
			}
		})
	}
}

func TestWriteToReaderWriterToError(t *testing.T) {
	want := errors.New("boom")
	wr := newWriteToReader(&errWriterTo{err: want})
	defer wr.Close()

	if err := wr.Join(); err != want {
		t.Errorf("Join, got %v, want %v", err, want)
	}
}

func TestWriteToReaderCloseBeforeRead(t *testing.T) {
	// Payload large enough that WriteTo must block waiting for a reader.
	big := bytes.Repeat([]byte("x"), 1<<20)
	wr := newWriteToReader(bytes.NewReader(big))

	if err := wr.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	if err := wr.Join(); err != io.ErrClosedPipe {
		t.Errorf("Join, got %v, want %v", err, io.ErrClosedPipe)
	}
}
