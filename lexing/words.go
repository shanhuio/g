package lexing

import (
	"io"
)

// Token types for the example lexer.
const (
	Word = iota
	Punc
)

func lexWord(x *Lexer) *Token {
	r := x.Rune()
	if IsLetter(r) || IsDigit(r) {
		// it is a word
		for {
			r, _ := x.Next()
			if x.Ended() || !(IsLetter(r) || IsDigit(r)) {
				break
			}
		}
		return x.MakeToken(Word)
	}

	x.Next()
	return x.MakeToken(Punc)
}

// NewWordLexer returns an example lexer that parses a file
// into words and punctuations.
func NewWordLexer(file string, r io.Reader) *Lexer {
	ret := MakeLexer(file, r, lexWord)
	ret.IsWhite = IsWhiteOrEndl
	return ret
}
