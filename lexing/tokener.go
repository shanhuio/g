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

// Tokener is token emitting interface.
type Tokener interface {
	// Token returns the next token
	Token() *Token

	// Errs returns the error list on tokening
	Errs() []*Error
}

// NewTokener creates a new tokener from LexFunc x and WhiteFunc w.
func NewTokener(f string, r io.Reader, x LexFunc, w WhiteFunc) Tokener {
	ret := NewLexer(f, r)
	ret.LexFunc = x
	ret.IsWhite = w

	return ret
}

// Tokens takes a lexer that is already setup and returns
// its tokens and errors.
func Tokens(tokener Tokener) ([]*Token, []*Error) {
	var ret []*Token
	for {
		t := tokener.Token()
		ret = append(ret, t)
		if t.Type == EOF {
			break
		}
	}
	return ret, tokener.Errs()
}
