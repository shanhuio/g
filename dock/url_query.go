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

package dock

import (
	"net/url"
	"path"
)

func urlQuery(p string, q url.Values) string {
	u := &url.URL{Path: p}
	if len(q) != 0 {
		u.RawQuery = q.Encode()
	}
	return u.String()
}

func apiURLQuery(p string, q url.Values) string {
	return urlQuery(path.Join(apiVersion, p), q)
}

func singleQuery(k, v string) url.Values {
	q := make(url.Values)
	q.Add(k, v)
	return q
}
