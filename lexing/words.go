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
	"io"
)

// Token types for the example lexer.
const (
	Word = iota
	Punc
)

func lexWord(x *Lexer) *Token {
	r := x.Rune()
	if IsLetter(r) || IsDigit(r) {
		// it is a word
		for {
			r, _ := x.Next()
			if x.Ended() || !(IsLetter(r) || IsDigit(r)) {
				break
			}
		}
		return x.MakeToken(Word)
	}

	x.Next()
	return x.MakeToken(Punc)
}

// NewWordLexer returns an example lexer that parses a file
// into words and punctuations.
func NewWordLexer(file string, r io.Reader) *Lexer {
	ret := MakeLexer(file, r, lexWord)
	ret.IsWhite = IsWhiteOrEndl
	return ret
}
