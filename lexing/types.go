package lexing

import (
	"fmt"
)

// Types is a registrar of token type names
type Types struct {
	names map[int]string
}

// NewTypes makes a new registrar. It will auto register the default tokens.
func NewTypes() *Types {
	ret := new(Types)
	ret.names = make(map[int]string)

	ret.Register(EOF, "eof")
	ret.Register(Comment, "comment")
	ret.Register(Illegal, "illegal")

	return ret
}

// Register registers a type with a name.
// If the type is already registered, it panics.
func (types *Types) Register(t int, name string) {
	if _, found := types.names[t]; found {
		panic("token already registered")
	}

	types.names[t] = name
}

// Name resolves the name of a type.
func (types *Types) Name(t int) string {
	if ret, found := types.names[t]; found {
		return ret
	}

	return fmt.Sprintf("<T%d>", t)
}
