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

package oauth2

import (
	"net/url"
	"strings"

	"shanhu.io/g/aries"
	"shanhu.io/g/errcode"
)

// ParseRedirect parses an in-site redirection URL.
// The server parts (scheme, host, port, user info) are discarded.
func ParseRedirect(redirect string) (string, error) {
	if redirect == "" {
		return "", nil
	}

	u, err := url.Parse(redirect)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(u.Path, "/") {
		return "", errcode.InvalidArgf(
			"redirect path part %q is not absolute", u.Path,
		)
	}

	return aries.DiscardURLServerParts(u).String(), nil
}
