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
