package lexing

import (
	"bufio"
	"io"
)

// runeScanner is a rune scanner that scans runes from a file,
// and at the same time tracks the reading position.
type runeScanner struct {
	file string
	line int
	col  int

	r *bufio.Reader

	Err  error // any error encountered
	Rune rune  // the rune just read

	closed bool
}

// newRuneScanner creates a scanner.
func newRuneScanner(file string, r io.Reader) *runeScanner {
	return &runeScanner{
		file: file,
		r:    bufio.NewReader(r),
		line: 1,
		col:  0,
	}
}

// scan reads in the next rune to s.Rune.  It closes the reader automatically
// when it reaches the end of file or when an error occurs.
func (s *runeScanner) scan() bool {
	if s.closed {
		panic("scanning on closed rune scanner")
	}

	wasEndline := s.Rune == '\n'

	s.Rune, _, s.Err = s.r.ReadRune()

	if s.Err != nil {
		s.closed = true
		return false
	}

	if wasEndline {
		s.line++
		s.col = 1
	} else {
		s.col++
	}

	return true
}

// pos returns the current position in the file.
func (s *runeScanner) pos() *Pos { return &Pos{s.file, s.line, s.col} }
