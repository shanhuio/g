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

package states

import (
	"encoding/json"

	"shanhu.io/pub/jsonx"
)

// GetJSON gets a JSON encoded state.
func GetJSON(ctx C, s States, key string, v interface{}) error {
	bs, err := s.Get(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, v)
}

// GetJSONX gets a JSONX encoded state.
func GetJSONX(ctx C, s States, key string, v interface{}) error {
	bs, err := s.Get(ctx, key)
	if err != nil {
		return err
	}
	return jsonx.Unmarshal(bs, v)
}

// PutJSON puts a JSON encoded state.
func PutJSON(ctx C, s States, key string, v interface{}) error {
	bs, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return s.Put(ctx, key, bs)
}
