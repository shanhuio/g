package signer

import (
	"testing"

	"time"
)

func TestTimeSigner(t *testing.T) {
	s := NewTimeSigner(nil, time.Second*5)
	now := time.Unix(0, 0)
	s.TimeFunc = func() time.Time { return now }

	token := s.Token()
	if !s.Check(token) {
		t.Errorf("token should be valid")
	}

	now = time.Unix(1, 0)
	if !s.Check(token) {
		t.Errorf("token should be still valid")
	}

	now = time.Unix(10, 0)
	if s.Check(token) {
		t.Errorf("token should be invalid")
	}

	now = time.Unix(-10, 0)
	if s.Check(token) {
		t.Errorf("token should be invalid")
	}
}
