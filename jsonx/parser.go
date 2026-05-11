package jsonx

import (
	"fmt"
	"io"

	"shanhu.io/g/lexing"
)

type parser struct {
	f string
	x lexing.Tokener
	*lexing.Parser
}

func newParser(f string, r io.Reader) (*parser, *lexing.Recorder) {
	p := &parser{f: f}
	x := tokener(f, r)
	rec := lexing.NewRecorder(x)
	p.x = lexing.NewCommentRemover(rec)
	p.Parser = lexing.NewParser(p.x, tokTypes)
	return p, rec
}

func (p *parser) seeOp(ops ...string) bool {
	t := p.Token()
	if t.Type != tokOperator {
		return false
	}
	for _, op := range ops {
		if t.Lit == op {
			return true
		}
	}
	return false
}

func tokenTypeStr(t *lexing.Token) string {
	if t.Type == tokOperator {
		return fmt.Sprintf("'%s'", t.Lit)
	}
	return typeStr(t.Type)
}

func (p *parser) expectOp(op string) *lexing.Token {
	if p.InError() {
		return nil
	}
	t := p.Token()
	if t.Type != tokOperator || t.Lit != op {
		p.CodeErrorfHere(
			"jsonx.expectOp", "expect '%s', got %s", op, tokenTypeStr(t),
		)
		return nil
	}
	return p.Shift()
}
