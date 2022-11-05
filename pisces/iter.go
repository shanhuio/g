// Copyright (C) 2022  Shanhu Tech Inc.
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
