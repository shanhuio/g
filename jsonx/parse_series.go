package jsonx

import (
	"shanhu.io/g/lexing"
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
