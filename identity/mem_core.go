// Copyright (C) 2023  Shanhu Tech Inc.
//
// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU Affero General Public License as published by the
// Free Software Foundation, either version 3 of the License, or (at your
// option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.  See the GNU Affero General Public License
// for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package identity

import (
	"encoding/json"
	"time"

	"shanhu.io/g/errcode"
)

type memStore struct {
	bs []byte
}

func (s *memStore) Check() (bool, error) {
	return s.bs != nil, nil
}

func (s *memStore) Save(v interface{}) error {
	bs, err := json.Marshal(v)
	if err != nil {
		return errcode.Annotate(err, "marshal")
	}
	s.bs = bs
	return nil
}

func (s *memStore) Load(v interface{}) error {
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
