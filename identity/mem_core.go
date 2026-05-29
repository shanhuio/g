package identity

import (
	"encoding/json"
	"time"

	"shanhu.io/std/errcode"
)

type memStore struct {
	bs []byte
}

func (s *memStore) Check() (bool, error) {
	return s.bs != nil, nil
}

func (s *memStore) Save(v any) error {
	bs, err := json.Marshal(v)
	if err != nil {
		return errcode.Annotate(err, "marshal")
	}
	s.bs = bs
	return nil
}

func (s *memStore) Load(v any) error {
	if len(s.bs) == 0 {
		return errcode.NotFoundf("identity not initialized")
	}
	if err := json.Unmarshal(s.bs, v); err != nil {
		return errcode.Annotate(err, "unmarshal")
	}
	return nil
}

// NewMemCore creates a new simple core that saves states in memory. It is
// useful for temporary testing.
func NewMemCore(t func() time.Time) Core {
	return NewSimpleCore(new(memStore), t)
}
