package lexing

import (
	"fmt"
)

// Pos is the file line position in a file
type Pos struct {
	File string
	Line int
	Col  int
}

func (p *Pos) String() string {
	if p.File == "" {
		return fmt.Sprintf("%d:%d", p.Line, p.Col)
	}
	return fmt.Sprintf("%s:%d:%d", p.File, p.Line, p.Col)
}
