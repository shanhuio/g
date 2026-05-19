package signer

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"time"
)

// Sessions signs a session data so that the server can run statelessly.
type Sessions struct {
	s          *Signer
	ttl        time.Duration
	refreshTTL time.Duration

	// TimeFunc is an optional function for reading the current timestamp.
	// When it is nil, the Sessions object uses time.Now().
	TimeFunc func() time.Time
}

func refreshTTL(ttl time.Duration) time.Duration {
	if ttl <= 0 {
		return time.Duration(0)
	}
	return ttl / 5
}

// NewSessions creates a new session store.
func NewSessions(key []byte, ttl time.Duration) *Sessions {
	return &Sessions{
		s:          New(key),
		ttl:        ttl,
		refreshTTL: refreshTTL(ttl),
	}
}

// New creates a new session with some data.
func (s *Sessions) New(data []byte, ttl time.Duration) (string, time.Time) {
	buf := new(bytes.Buffer)

	if ttl <= 0 || ttl > s.ttl {
		ttl = s.ttl
	}

	// write the timestamp
	expires := now(s.TimeFunc).Add(ttl)
	ts := make([]byte, timestampLen)
	binary.LittleEndian.PutUint64(ts, uint64(expires.UnixNano()))
	buf.Write(ts)

	if data != nil {
		buf.Write(data)
	}

	return s.s.SignHex(buf.Bytes()), expires
}

// NewJSON creates a new session with a JSON marshallabe data.
func (s *Sessions) NewJSON(data any) (string, time.Time, error) {
	bs, err := json.Marshal(data)
	if err != nil {
		var t time.Time
		return "", t, err
	}

	ret, expires := s.New(bs, 0)
	return ret, expires, nil
}

// Check checks if it is a signed data
func (s *Sessions) Check(session string) ([]byte, time.Duration, bool) {
	ok, bs := s.s.CheckHex(session)
	if !ok {
		return nil, 0, false
	}

	if len(bs) < timestampLen {
		return nil, 0, false
	}

	ts := int64(binary.LittleEndian.Uint64(bs[:timestampLen]))
	expire := time.Unix(0, ts)
	timeNow := now(s.TimeFunc)

	if !timeNow.Before(expire) {
		return nil, 0, false
	}

	return bs[timestampLen:], expire.Sub(timeNow), true
}

// CheckJSON checks if the session is valid and unmarshals if it is.
// It will return false if it is fails to unmarshal.
func (s *Sessions) CheckJSON(session string, dat any) bool {
	bs, _, ok := s.Check(session)
	if !ok {
		return false
	}
	return json.Unmarshal(bs, dat) == nil
}

// NewState creates a new state, which is a session with no data.
func (s *Sessions) NewState() string {
	ret, _ := s.New(nil, 0)
	return ret
}

// CheckState checks if it is a signed session with no data.
func (s *Sessions) CheckState(session string) bool {
	bs, _, ok := s.Check(session)
	if !ok {
		return false
	}
	return len(bs) == 0
}

// NeedRefresh returns if it is recommended to refresh for a new token.
func (s *Sessions) NeedRefresh(ttl time.Duration) bool {
	return ttl < s.refreshTTL
}
