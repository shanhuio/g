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

// WhiteFunc is a function type that checks if a rune is white space.
type WhiteFunc func(r rune) bool

// IsWhite is the default IsWhite function for a lexer. Returns true on spaces,
// \t and \r.  Returns false on \n.
func IsWhite(r rune) bool {
	return r == ' ' || r == '\t' || r == '\r'
}

// IsWhiteOrEndl is another IsWhite function that also returns true for \n.
func IsWhiteOrEndl(r rune) bool {
	return IsWhite(r) || r == '\n'
}
