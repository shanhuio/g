package lexing

import (
	"strings"
	"testing"
)

func TestWordLexer(t *testing.T) {
	x := NewWordLexer("t1.txt", strings.NewReader("hello, world!"))
	toks, errs := Tokens(x)
	if len(errs) != 0 {
		t.Errorf("unexpected errors: %v", errs)
		return
	}
	if len(toks) != 5 {
		t.Errorf("want 5 tokens, got %d", len(toks))
		return
	}

	for i, s := range []string{"hello", ",", "world", "!", ""} {
		lit := toks[i].Lit
		if s != lit {
			t.Errorf("token %d want %q, got %q", i, s, lit)
		}
	}

	for i, want := range []int{Word, Punc, Word, Punc, EOF} {
		typ := toks[i].Type
		if want != typ {
			t.Errorf("token %d want type %d, got %d", i, want, typ)
		}
	}

	x = NewWordLexer("t2.txt", strings.NewReader("123	#a $%\n\rd^"))
	toks, errs = Tokens(x)
	if len(errs) != 0 {
		t.Errorf("unexpected errors: %v", errs)
		return
	}
	if len(toks) != 8 {
		t.Errorf("want 8 tokens, got %d", len(toks))
		return
	}
	for i, s := range []string{"123", "#", "a", "$", "%", "d", "^", ""} {
		lit := toks[i].Lit
		if s != lit {
			t.Errorf("token %d want %q, got %q", i, s, lit)
		}
	}

	for i, want := range []int{Word, Punc, Word, Punc, Punc, Word, Punc, EOF} {
		typ := toks[i].Type
		if want != typ {
			t.Errorf("token %d want type %d, got %d", i, want, typ)
		}
	}
}
