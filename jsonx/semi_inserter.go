package jsonx

import (
	"shanhu.io/g/lexing"
)

type semiInserter struct {
	x          lexing.Tokener
	save       *lexing.Token
	insertSemi bool
}

func newSemiInserter(x lexing.Tokener) *semiInserter {
	return &semiInserter{x: x}
}

func makeSemi(p *lexing.Pos, lit string) *lexing.Token {
	return &lexing.Token{Pos: p, Lit: lit, Type: tokSemi}
}

func (si *semiInserter) Token() *lexing.Token {
	if si.save != nil {
		t := si.save
		si.save = nil
		return t
	}

	for {
		t := si.x.Token()
		switch t.Type {
		case tokSemi:
			si.insertSemi = false
		case tokOperator:
			switch t.Lit {
			case "}", "]":
				si.insertSemi = true
			default:
				si.insertSemi = false
			}
		case lexing.EOF:
			if si.insertSemi {
				si.insertSemi = false
				si.save = t
				return makeSemi(t.Pos, "")
			}
		case tokEndl:
			if si.insertSemi {
				si.insertSemi = false
				return makeSemi(t.Pos, "\n")
			}
			continue
		case lexing.Comment:
			// do nothing
		default:
			si.insertSemi = true
		}

		return t
	}
}

func (si *semiInserter) Errs() []*lexing.Error {
	return si.x.Errs()
}
