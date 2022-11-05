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

package fmtutil

import (
	"bytes"
	"fmt"
	"reflect"
)

// Join joins a slice of stuff into a string with sep as the
// separator.
func Join(slice interface{}, sep string) string {
	t := reflect.TypeOf(slice)
	if t.Kind() != reflect.Slice {
		return fmt.Sprint(slice)
	}

	v := reflect.ValueOf(slice)
	n := v.Len()
	if n == 0 {
		return ""
	} else if n == 1 {
		return fmt.Sprint(v.Index(0).Interface())
	}

	buf := new(bytes.Buffer)
	for i := 0; i < n; i++ {
		x := v.Index(i).Interface()
		if i > 0 {
			fmt.Fprint(buf, sep)
		}
		fmt.Fprint(buf, x)
	}
	return buf.String()
}
