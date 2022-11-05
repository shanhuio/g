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

// LexNumber lexes a number usign golang's number format.
func LexNumber(x *Lexer, tokInt, tokFloat int) *Token {
	isFloat := false
	start := x.Rune()
	if !IsDigit(start) {
		panic("not starting with a number")
	}

	x.Next()
	r := x.Rune()
	if start == '0' && r == 'x' {
		x.Next()
		for IsHexDigit(x.Rune()) {
			x.Next()
		}
	} else {
		for IsDigit(x.Rune()) {
			x.Next()
		}
		if x.Rune() == '.' {
			isFloat = true
			x.Next()
			for IsDigit(x.Rune()) {
				x.Next()
			}
		}
		if x.Rune() == 'e' || x.Rune() == 'E' {
			isFloat = true
			x.Next()
			if IsDigit(x.Rune()) || x.Rune() == '-' {
				x.Next()
			}
			for IsDigit(x.Rune()) {
				x.Next()
			}
		}
	}
	if isFloat {
		return x.MakeToken(tokFloat)
	}
	return x.MakeToken(tokInt)
}
