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

package keyreg

import (
	"net/url"

	"shanhu.io/pub/errcode"
)

// OpenKeyRegistry connects to a keystore based on the given URL string.
func OpenKeyRegistry(urlStr string) (KeyRegistry, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "http", "https":
		u, err := url.Parse(urlStr)
		if err != nil {
			return nil, err
		}
		return NewWebKeyRegistry(u), nil
	case "file", "":
		r, err := NewDirKeyRegistry(u.Path)
		if err != nil {
			return nil, err
		}
		return r, nil
	}
	return nil, errcode.InvalidArgf("unsupported url scheme: %q", u.Scheme)
}
