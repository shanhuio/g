package signer

import (
	"testing"

	"time"
)

func TestStates(t *testing.T) {
	s := NewSessions(nil, time.Second)
	state := s.NewState()
	if !s.CheckState(state) {
		t.Errorf("check on state %q failed", state)
	}

	if s.CheckState("") {
		t.Errorf("check on empty state is passing")
	}
}

func TestStatesExpire(t *testing.T) {
	const ttl = time.Second
	s := NewSessions(nil, ttl)
	now := time.Unix(0, 0)
	s.TimeFunc = func() time.Time { return now }
	state := s.NewState()
	t.Log("state: ", state)

	now = now.Add(ttl).Add(-time.Nanosecond)
	if !s.CheckState(state) {
		t.Errorf("check on state %q failed", state)
	}

	now = now.Add(time.Nanosecond)
	if s.CheckState(state) {
		t.Errorf("check passed, should fail because of time out")
	}
}
