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

// KeywordSet creates a keyword set.
func KeywordSet(words ...string) map[string]struct{} {
	ret := make(map[string]struct{})
	for _, k := range words {
		ret[k] = struct{}{}
	}
	return ret
}

// Keyworder contains idents into keywords
type Keyworder struct {
	tokener Tokener

	Keywords map[string]struct{}
	Ident    int
	Keyword  int
}

// NewKeyworder creates a new tokener that changes the type
// of a token into keywords if it is in the keyword map.
func NewKeyworder(tok Tokener) *Keyworder {
	return &Keyworder{tokener: tok}
}

// Token returns the next token, while replacing ident types into
// keyword types if the token is in the keyword set.
func (kw *Keyworder) Token() *Token {
	ret := kw.tokener.Token()
	if kw.Keywords != nil && ret.Type == kw.Ident {
		_, ok := kw.Keywords[ret.Lit]
		if ok {
			ret.Type = kw.Keyword
		}
	}

	return ret
}

// Errs returns the error list on tokening.
func (kw *Keyworder) Errs() []*Error {
	return kw.tokener.Errs()
}
