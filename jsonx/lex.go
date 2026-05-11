package jsonx

import (
	"io"

	"shanhu.io/g/lexing"
)

func lexOperator(x *lexing.Lexer, r rune) *lexing.Token {
	switch r {
	case '{', '}', '[', ']', ',', ':', '+', '-', '.':
		/* do nothing */
	case '/':
		r2 := x.Rune()
		if r2 == '/' || r2 == '*' {
			return lexing.LexComment(x)
		}
	case ';':
		return x.MakeToken(tokSemi)
	default:
		return nil
	}
	return x.MakeToken(tokOperator)
}

func lexJSONX(x *lexing.Lexer) *lexing.Token {
	r := x.Rune()
	if x.IsWhite(r) {
		panic("incorrect token start")
	}

	switch r {
	case '\n':
		x.Next()
		return x.MakeToken(tokEndl)
	case '"':
		return lexing.LexString(x, tokString, '"')
	case '`':
		return lexing.LexRawString(x, tokString)
	}

	if lexing.IsDigit(r) {
		return lexing.LexNumber(x, tokInt, tokFloat)
	}
	if lexing.IsIdentLetter(r) {
		return lexing.LexIdent(x, tokIdent)
	}

	x.Next()
	t := lexOperator(x, r)
	if t != nil {
		return t
	}

	x.CodeErrorf("jsonx.illegalChar", "illegal char %q", r)
	return x.MakeToken(lexing.Illegal)
}

var keywords = lexing.KeywordSet("true", "false", "null")

func tokener(f string, r io.Reader) lexing.Tokener {
	x := lexing.MakeLexer(f, r, lexJSONX)
	si := newSemiInserter(x)
	kw := lexing.NewKeyworder(si)
	kw.Ident = tokIdent
	kw.Keyword = tokKeyword
	kw.Keywords = keywords
	return kw
}
