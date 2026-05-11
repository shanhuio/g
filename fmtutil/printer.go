package fmtutil

import (
	"bytes"
	"io"
)

// Printer is a write filter that supports shift tabbing
type Printer struct {
	out     io.Writer
	e       error
	midLine bool

	indent    int
	indentStr string
}

// NewPrinter creates a new printer
func NewPrinter(out io.Writer) *Printer {
	ret := new(Printer)
	ret.out = out
	ret.indentStr = "    "
	return ret
}

func (p *Printer) write(buf []byte) {
	if p.e != nil {
		return
	}
	_, p.e = p.out.Write(buf)
}

func (p *Printer) writeBytes(buf []byte) {
	if len(buf) == 0 {
		return
	}

	if !p.midLine {
		for j := 0; j < p.indent; j++ {
			p.write([]byte(p.indentStr))
		}
	}

	p.midLine = true
	p.write(buf)
}

func (p *Printer) writeEndl() {
	p.write([]byte("\n"))
	p.midLine = false
}

// Write writes the buf. It adds indent before each line.
func (p *Printer) Write(buf []byte) (int, error) {
	lines := bytes.Split(buf, []byte("\n"))

	for i, line := range lines {
		if i > 0 {
			p.writeEndl()
		}

		p.writeBytes(line)
	}

	return len(buf), nil
}

// Tab indents in one level
func (p *Printer) Tab() { p.indent++ }

// ShiftTab indents out one level
func (p *Printer) ShiftTab() {
	if p.indent > 0 {
		p.indent--
	}
}

// Err returns the first error on printing.
func (p *Printer) Err() error { return p.e }
