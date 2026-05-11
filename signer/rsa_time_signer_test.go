package signer

import (
	"testing"

	"crypto/rand"
	"crypto/rsa"
	"time"
)

func TestRSATimeSigner(t *testing.T) {
	const size = 2048
	key, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()

	s := NewRSATimeSigner(&key.PublicKey, time.Second)
	clock := now
	s.TimeFunc = func() time.Time { return clock }

	b, err := rsaSignTime(key, now)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		diff time.Duration
		ok   bool
	}{
		{0, true},
		{time.Second / 2, true},
		{-time.Second / 2, true},
		{time.Second * 2, false},
		{-time.Second * 2, false},
	} {
		clock = now.Add(test.diff)
		err := s.Check(b)
		if err != nil && test.ok {
			t.Errorf("unexpected error for time diff %s: %s", test.diff, err)
		} else if err == nil && !test.ok {
			t.Errorf("timestamp should be out of window for diff %s", test.diff)
		}
	}

	clock = now
	anotherKey, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		t.Fatal(err)
	}
	b, err = rsaSignTime(anotherKey, now)
	if err != nil {
		t.Fatal(err)
	}
	if s.Check(b) == nil {
		t.Errorf("signer should not valid")
	}
}
