package pisces

import (
	"encoding/json"
)

// Iter is an interator.
type Iter struct {
	Make func() interface{}
	Do   func(cls string, v interface{}) error
}

// KVPartial specifies a part of a query result.
type KVPartial struct {
	Offset uint64
	N      uint64
	Desc   bool
}

func (it *Iter) doWalk(_, cls string, bs []byte) error {
	v := it.Make()
	if err := json.Unmarshal(bs, v); err != nil {
		return err
	}
	return it.Do(cls, v)
}
