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

package textbox

import (
	"bufio"
	"io"

	"shanhu.io/pub/lexing"
)

// TabSize is the indent size for each tab
const TabSize = 4

func runeWidth(r rune) int {
	switch r {
	case '\t':
		return TabSize
	case '\n', '\r':
		return 0
	}
	return 1
}

// Rect returns the display of a text line.
// Ends of lines are ignored.
func Rect(r io.Reader) (nline, maxWidth int, err error) {
	br := bufio.NewReader(r)
	nline = 0
	curWidth := 0
	maxWidth = 0

	for {
		r, _, err := br.ReadRune()
		if err == io.EOF {
			if curWidth > 0 {
				nline++
			}
			break
		} else if err != nil {
			return 0, 0, err
		}

		if r == '\n' {
			nline++
			if curWidth > maxWidth {
				maxWidth = curWidth
			}
			curWidth = 0
		} else {
			curWidth += runeWidth(r)
		}
	}

	if curWidth > maxWidth {
		maxWidth = curWidth
	}

	return nline, maxWidth, nil
}

// CheckRect checks if a program is within a rectangular area.
func CheckRect(file string, r io.Reader, h, w int) []*lexing.Error {
	br := bufio.NewReader(r)
	line := 0
	col := 0

	errs := lexing.NewErrorList()

	pos := func() *lexing.Pos {
		return &lexing.Pos{
			File: file,
			Line: line + 1,
			Col:  col + 1,
		}
	}
	newLine := func() {
		if col > w {
			errs.Errorf(
				pos(),
				"this line is too wide. it has %d chars; the limit is %d",
				col, w,
			)
		}
		line++
		col = 0
	}

	for {
		r, _, e := br.ReadRune()
		if e == io.EOF {
			if col > 0 {
				newLine()
			}
			break
		} else if lexing.LogError(errs, e) {
			break
		}

		if r == '\n' {
			newLine()
		} else {
			col += runeWidth(r)
		}
	}

	if line > h {
		errs.Errorf(pos(), "has too many lines; the limit is %d", h)
	}

	return errs.Errs()
}
