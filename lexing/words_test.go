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
	"strings"
	"testing"
)

func TestWordLexer(t *testing.T) {
	x := NewWordLexer("t1.txt", strings.NewReader("hello, world!"))
	toks, errs := Tokens(x)
	if len(errs) != 0 {
		t.Errorf("unexpected errors: %v", errs)
		return
	}
	if len(toks) != 5 {
		t.Errorf("want 5 tokens, got %d", len(toks))
		return
	}

	for i, s := range []string{"hello", ",", "world", "!", ""} {
		lit := toks[i].Lit
		if s != lit {
			t.Errorf("token %d want %q, got %q", i, s, lit)
		}
	}

	for i, want := range []int{Word, Punc, Word, Punc, EOF} {
		typ := toks[i].Type
		if want != typ {
			t.Errorf("token %d want type %d, got %d", i, want, typ)
		}
	}

	x = NewWordLexer("t2.txt", strings.NewReader("123	#a $%\n\rd^"))
	toks, errs = Tokens(x)
	if len(errs) != 0 {
		t.Errorf("unexpected errors: %v", errs)
		return
	}
	if len(toks) != 8 {
		t.Errorf("want 8 tokens, got %d", len(toks))
		return
	}
	for i, s := range []string{"123", "#", "a", "$", "%", "d", "^", ""} {
		lit := toks[i].Lit
		if s != lit {
			t.Errorf("token %d want %q, got %q", i, s, lit)
		}
	}

	for i, want := range []int{Word, Punc, Word, Punc, Punc, Word, Punc, EOF} {
		typ := toks[i].Type
		if want != typ {
			t.Errorf("token %d want type %d, got %d", i, want, typ)
		}
	}
}
