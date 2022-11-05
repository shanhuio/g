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

// Remover removes a particular type of token from a token stream
type Remover struct {
	Tokener
	t int
}

// NewRemover creates a new remover that removes token of type t
func NewRemover(t Tokener, typ int) *Remover {
	if typ == EOF {
		panic("cannot remove EOF")
	}

	ret := new(Remover)
	ret.Tokener = t
	ret.t = typ

	return ret
}

// Token implements the Tokener interface but only returns
// token that is not the particular type.
func (r *Remover) Token() *Token {
	for {
		ret := r.Tokener.Token()
		if ret.Type != r.t {
			return ret
		}
	}
}

// NewCommentRemover creates a new remover that removes token
func NewCommentRemover(t Tokener) *Remover {
	return NewRemover(t, Comment)
}
