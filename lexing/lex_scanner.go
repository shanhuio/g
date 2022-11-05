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

package lexing

import (
	"bytes"
	"io"
)

// lexScanner parses a file input stream into tokens.
type lexScanner struct {
	s     *runeScanner
	valid bool

	pos *Pos
	buf *bytes.Buffer
}

// newLexScanner creates a new lexer.
func newLexScanner(file string, r io.Reader) *lexScanner {
	ret := &lexScanner{
		s:   newRuneScanner(file, r),
		buf: new(bytes.Buffer),
	}
	ret.pos = ret.s.pos()
	return ret
}

// next pushes the current rune (if valid) into the buffer,
// and returns the next rune or error from scanning the input
// stream.
func (s *lexScanner) next() (rune, error) {
	if s.valid {
		s.buf.WriteRune(s.s.Rune) // push into the buffer
		s.valid = false
	}

	if !s.s.scan() {
		return 0, s.s.Err
	}

	s.valid = true
	return s.s.Rune, nil
}

// accept returns the string buffered, and the starting position
// of the string.
func (s *lexScanner) accept() (string, *Pos) {
	ret := s.buf.String()
	s.buf.Reset()
	pos := s.pos

	s.pos = s.s.pos()
	return ret, pos
}

// buffered returns the current buffered string in the
// scanner
func (s *lexScanner) buffered() string {
	return s.buf.String()
}

// startPos returns the position of the buffer start.
func (s *lexScanner) startPos() *Pos { return s.pos }
