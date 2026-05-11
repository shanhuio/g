// Package strtoken parses a line into string tokens.  Tokens are separated by
// white spaces, and can be double quoted for tokens that contains white
// spaces.
package strtoken

import (
	"io"
	"strconv"
	"strings"

	"shanhu.io/g/lexing"
)

func lexShell(x *lexing.Lexer) *lexing.Token {
	r := x.Rune()
	if x.IsWhite(r) {
		panic("incorrect token start")
	}

	switch r {
	case '"':
		return lexing.LexString(x, str, '"')
	}

	if isBareRune(r) {
		return lexBare(x)
	}

	x.Errorf("illegal char %q", r)
	x.Next()
	return x.MakeToken(lexing.Illegal)
}

func newLexer(file string, r io.Reader) *lexing.Lexer {
	return lexing.MakeLexer(file, r, lexShell)
}

// Parse parses a line into a sequence of tokens.
// A token
func Parse(line string) ([]string, []*lexing.Error) {
	r := strings.NewReader(line)
	x := newLexer("", r)
	toks := lexing.TokenAll(x)
	if errs := x.Errs(); errs != nil {
		return nil, errs
	}

	var ret []string
	var errs []*lexing.Error
	for _, t := range toks {
		switch t.Type {
		case bare:
			ret = append(ret, t.Lit)
		case str:
			v, err := strconv.Unquote(t.Lit)
			if err != nil {
				errs = append(errs, &lexing.Error{
					Pos:  t.Pos,
					Err:  err,
					Code: "shellarg.invalidStr",
				})
			} else {
				ret = append(ret, v)
			}
		}
	}

	if errs != nil {
		return nil, errs
	}
	return ret, nil
}
