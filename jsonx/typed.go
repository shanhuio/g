package jsonx

import (
	"shanhu.io/g/lexing"
)

// TypeMaker is a function that makes an interface based on the given type.
type TypeMaker func(t string) interface{}

// Typed is an item in a typed list.
type Typed struct {
	Type string
	V    interface{}
	Pos  *lexing.Pos
}
