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
	"testing"

	"strings"
)

func TestLexComment(t *testing.T) {
	for _, s := range []string{
		"",
		"     ",
		"// abc",
		"/*abc*/",
		"   // abc",
		"   /* abc \n abc */",
	} {
		t.Logf("parsing %q", s)
		x := NewCommentLexer("a.txt", strings.NewReader(s))
		toks, errs := Tokens(x)
		if len(errs) != 0 {
			t.Errorf("unxpected errors: %v", errs)
		}
		if len(toks) > 2 {
			t.Errorf("want 0 or 1 token, got %d", len(toks)-1)
		}
	}

	for _, s := range []string{
		"/*//*a/",    // bad end
		"/*\n\n",     // bad end
		"//abc\nabc", // cross line
		"/*abc*/abc", // closed
		"/abc",       // invalid start
		"/",          //
	} {
		t.Logf("parsing %q", s)
		x := NewCommentLexer("a.txt", strings.NewReader(s))
		_, errs := Tokens(x)
		if len(errs) == 0 {
			t.Errorf("expect error, but got nothing")
		}
	}
}
