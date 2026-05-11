package timeutil

import (
	"testing"

	"crypto/rand"
	"time"
)

func TestChallenge(t *testing.T) {
	now := time.Now()

	r := rand.Reader
	ch, err := NewChallenge(now, r)
	if err != nil {
		t.Fatal(err)
	}

	got := Time(ch.T)
	if !now.Equal(got) {
		t.Errorf("got timestamp %q, want %q", got, now)
	}
	if ch.N == "" {
		t.Errorf("nounce is empty")
	}

	ch2, err := NewChallenge(now, r)
	if err != nil {
		t.Fatal("get second challenge: ", err)
	}

	if ch2.N == ch.N {
		t.Errorf("got same nounce: %q", ch.N)
	}
}
