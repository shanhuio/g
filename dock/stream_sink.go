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

package dock

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"shanhu.io/pub/errcode"
)

type progressDetail struct {
	Current int64 `json:"current"`
	Total   int64 `json:"total"`
}

type streamMessage struct {
	From           string          `json:"from,omitempty"`
	ID             string          `json:"id,omitempty"`
	Stream         string          `json:"stream,omitempty"`
	ProgressDetail *progressDetail `json:"progressDetail,omitempty"`

	Status string `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}

type streamSink struct {
	io.WriteCloser
	r    *io.PipeReader
	w    *io.PipeWriter
	done chan error
	out  io.Writer
}

func newStreamSink() *streamSink {
	r, w := io.Pipe()

	ret := &streamSink{
		WriteCloser: w,
		r:           r,
		w:           w,
		done:        make(chan error, 1),
	}
	go func() { ret.done <- printStreamMessage(r, os.Stderr) }()
	return ret
}

type progressPrinter struct {
	ids []string
	out io.Writer
}

func (p *progressPrinter) Print(m *streamMessage) {
	if m.Stream != "" {
		fmt.Fprint(p.out, m.Stream)
		p.ids = nil // reset progress printer
	}
	if m.ID == "" {
		return
	}

	index := -1
	for i, id := range p.ids {
		if id == m.ID {
			index = i
			break
		}
	}
	if index == -1 {
		index = len(p.ids)
		p.ids = append(p.ids, m.ID)
		fmt.Fprintln(p.out)
	}

	n := len(p.ids)
	move := n - index
	// go up to line index
	fmt.Fprintf(p.out, "\033[%dA", move)
	fmt.Fprintf(p.out, "\r\033[K")

	if d := m.ProgressDetail; d != nil {
		fmt.Fprintf(p.out, "%s  %s", m.ID, m.Status)
		if d.Total > 0 {
			fmt.Fprintf(
				p.out, "  %.2f%% of %dB",
				float64(d.Current)/float64(d.Total)*100, d.Total,
			)
		}
	} else if m.Status != "" {
		fmt.Fprint(p.out, m.Status)
	}
	fmt.Fprintf(p.out, "\033[%dB\r", move)
}

func printStreamMessage(r io.ReadCloser, out io.Writer) error {
	defer r.Close()
	var streamErr error

	p := &progressPrinter{out: out}
	dec := json.NewDecoder(r)
	for dec.More() {
		m := new(streamMessage)
		if err := dec.Decode(m); err != nil {
			return err
		}
		if m.Error != "" {
			if streamErr != nil {
				log.Println(streamErr)
			}
			streamErr = errcode.Add(errcode.Internal, errors.New(m.Error))
		}
		p.Print(m)
	}
	return streamErr
}

func (s *streamSink) waitDone() error {
	if err := s.w.Close(); err != nil {
		return err
	}
	return <-s.done
}
