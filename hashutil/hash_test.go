package hashutil

import (
	"testing"

	"strings"
)

func TestHash(t *testing.T) {
	m := make(map[string]bool)
	addHash := func(h string) {
		if m[h] {
			t.Fatalf("hash conflict: %s", h)
		}
		m[h] = true
	}

	addHash(Hash(nil))
	addHash(HashStr("a"))
	addHash(HashStr("A"))
	addHash(HashStr("A "))
	addHash(HashStr("Hello"))

	const s = "something"
	h1 := HashStr(s)
	h2, err := HashReader(strings.NewReader(s))
	if err != nil {
		t.Fatal(err)
	}
	if h1 != h2 {
		t.Errorf("HashStr(%q) != HashReader(%q)", s, s)
	}
}

func TestHashFile(t *testing.T) {
	got, err := HashFile("testdata/testfile")
	if err != nil {
		t.Fatal(err)
	}

	want := HashStr("something\n")
	if want != got {
		t.Errorf("HashFile want %q, got %q", want, got)
	}
}
