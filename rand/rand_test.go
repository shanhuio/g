package rand

import (
	"testing"

	"bytes"
	"strings"
)

func TestBytes(t *testing.T) {
	bs1 := Bytes(8)
	bs2 := Bytes(8)

	if bytes.Equal(bs1, bs2) {
		t.Errorf("not so random: %v == %v", bs1, bs2)
	}
}

func TestLowerLetters(t *testing.T) {
	s1 := LowerLetters(16)
	s2 := LowerLetters(16)
	if s1 == s2 {
		t.Errorf("not so random: %q == %q", s1, s2)
	}
	if strings.ToLower(s1) != s1 {
		t.Errorf("contains non-lower case: %q", s1)
	}
}

func TestLetters(t *testing.T) {
	s1 := Letters(16)
	s2 := Letters(16)
	if s1 == s2 {
		t.Errorf("not so random: %q == %q", s1, s2)
	}
}

func TestDigits(t *testing.T) {
	s1 := Digits(10)
	s2 := Digits(10)
	if s1 == s2 {
		t.Errorf("not so random: %q == %q", s1, s2)
	}

	for _, s := range []string{s1, s2} {
		for _, r := range s {
			if !(r >= '0' && r <= '9') {
				t.Errorf("digits string %q contains %q", s, r)
			}
		}
	}
}
