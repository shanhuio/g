package strtoken

import (
	"shanhu.io/std/lexing"
)

func isBareRune(r rune) bool {
	if r == ' ' || r == '\n' || r == '\r' {
		return false
	}
	return true
}

func lexBare(x *lexing.Lexer) *lexing.Token {
	start := x.Rune()
	if !isBareRune(start) {
		panic("not starting with a bare rune")
	}

	for {
		x.Next()
		if x.Ended() || !isBareRune(x.Rune()) {
			break
		}
	}

	return x.MakeToken(bare)
}
