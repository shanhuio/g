package lexing

import (
	"fmt"
)

// Token defines a token structure.
type Token struct {
	Type int
	Lit  string
	Pos  *Pos
}

// Standard token types
const (
	EOF = -1 - iota
	Comment
	Illegal
)

func (t *Token) String() string {
	return fmt.Sprintf("'%s' (%v)", t.Lit, t.Pos)
}
