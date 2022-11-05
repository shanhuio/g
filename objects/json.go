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

package objects

import (
	"bytes"
	"encoding/json"
)

// ReadJSON reads an object blob and unmarshals it into a JSON object.
func ReadJSON(b Objects, k string, v interface{}) error {
	r, err := b.Open(k)
	if err != nil {
		return err
	}
	defer r.Close()
	if err := json.NewDecoder(r).Decode(v); err != nil {
		return err
	}
	return r.Close()
}

// CreateJSON creates an object blob that is the marshalling of the object.
func CreateJSON(b Objects, v interface{}) (string, error) {
	bs, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return b.Create(bytes.NewBuffer(bs))
}
