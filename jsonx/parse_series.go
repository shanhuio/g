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

package jsonx

import (
	"shanhu.io/pub/lexing"
)

func parseTypeName(p *parser) *typeName {
	if p.See(tokString) {
		tok := p.Shift()
		v := parseStringValue(p, tok)
		return &typeName{
			tok:  tok,
			name: v,
		}
	} else if p.See(tokIdent) {
		tok := p.Shift()
		return &typeName{
			tok:  tok,
			name: tok.Lit,
		}
	}

	t := p.Token()
	p.CodeErrorf(
		t.Pos, "jsonx.expectTypeName",
		"expect string or identifier, got %s", tokenTypeStr(t),
	)
	return nil
}

func parseSeries(p *parser) *series {
	s := new(series)

	for !p.See(lexing.EOF) {
		name := parseTypeName(p)
		if name == nil {
			p.SkipErrStmt(tokSemi)
			continue
		}

		value := parseValue(p)
		if p.SkipErrStmt(tokSemi) {
			continue
		}

		semi := p.Expect(tokSemi)
		s.entries = append(s.entries, &typedEntry{
			typ:   name,
			value: value,
			semi:  semi,
		})

		p.SkipErrStmt(tokSemi)
	}

	return s
}
