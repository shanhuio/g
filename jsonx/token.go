package jsonx

import (
	"shanhu.io/g/lexing"
)

const (
	tokKeyword = iota
	tokIdent
	tokString
	tokInt
	tokFloat
	tokOperator
	tokSemi
	tokEndl
)

var tokTypes = func() *lexing.Types {
	t := lexing.NewTypes()
	for _, e := range []struct {
		t    int
		name string
	}{
		{tokKeyword, "keyword"},
		{tokIdent, "identifier"},
		{tokString, "string"},
		{tokInt, "integer"},
		{tokFloat, "float"},
		{tokOperator, "operator"},
		{tokEndl, "end-line"},
		{tokSemi, "end-line"},
	} {
		t.Register(e.t, e.name)
	}
	return t
}()

func typeStr(t int) string { return tokTypes.Name(t) }
