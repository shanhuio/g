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
	"unicode"
)

func digitVal(r rune) int {
	switch {
	case '0' <= r && r <= '9':
		return int(r - '0')
	case 'a' <= r && r <= 'f':
		return int(r - 'a' + 10)
	case 'A' <= r && r <= 'F':
		return int(r - 'A' + 10)
	}
	return 16
}

func lexEscape(x *Lexer, quote rune) bool {
	var n int
	var base, max uint32
	if x.Ended() {
		x.Errorf("escape not terminated")
		return false
	}
	r := x.Rune()
	switch r {
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', quote:
		x.Next()
		return true
	case '0', '1', '2', '3', '4', '5', '6', '7':
		n, base, max = 3, 8, 255
	case 'x':
		x.Next()
		n, base, max = 2, 16, 255
	case 'u':
		x.Next()
		n, base, max = 4, 16, unicode.MaxRune
	case 'U':
		x.Next()
		n, base, max = 8, 16, unicode.MaxRune
	default:
		x.CodeErrorf("lexing.unknownESC", "unknown escape sequence")
		return false
	}

	var v uint32
	for i := 0; i < n; i++ {
		if x.Ended() {
			x.Errorf("escape not terminated")
			return false
		}

		r := x.Rune()
		d := uint32(digitVal(r))
		if d >= base {
			x.Errorf("illegal escape char %#U", r)
			return false
		}

		v *= base
		v += d

		x.Next()
	}

	if v > max || 0xD800 <= v && v < 0xE000 {
		x.Errorf("invalid unicode code point")
		return false
	}

	return true
}

// LexRawString parses a raw string token with type t, which is
// quoted in a pair of `
func LexRawString(x *Lexer, t int) *Token {
	if !x.See('`') {
		panic("incorrect raw string start")
	}

	x.Next()
	for {
		if x.Ended() {
			x.CodeErrorf("lexing.unexpectedEOF", "unexpected eof in raw string")
			break
		} else if x.See('`') {
			x.Next()
			break
		} else {
			x.Next()
		}
	}
	return x.MakeToken(t)
}

// LexString parses a string token with type t.
func LexString(x *Lexer, t int, q rune) *Token {
	if !(q == '\'' || q == '"') {
		panic("only support `'` or `\"`")
	} else if !x.See(q) {
		panic("incorrect string start")
	}

	n := 0
	x.Next()
	for {
		if x.Ended() {
			x.CodeErrorf("lexing.unexpectedEOF", "unexpected eof in string")
			break
		} else if x.See('\n') {
			x.CodeErrorf("lexing.unexpectedEndl", "unexpected endl in string")
			break
		} else if x.See(q) {
			x.Next()
			break
		} else if x.See('\\') {
			x.Next()
			lexEscape(x, q)
		} else {
			x.Next()
		}
		n++
	}

	if q == '\'' && n != 1 {
		x.CodeErrorf("lexing.illegalCharLit", "illegal char literal")
	}
	return x.MakeToken(t)
}
