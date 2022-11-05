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

package creds

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"strings"
)

func extractHost(s string) string {
	if strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https") {
		u, err := url.Parse(s)
		if err != nil {
			return s
		}
		return u.Host
	}
	return s
}

// Filename returns the creds file name for a particular domain.
func Filename(s string) string {
	s = extractHost(s)

	buf := new(bytes.Buffer)
	dashing := false

	out := func(r rune) {
		if dashing {
			buf.WriteString("-")
		}
		buf.WriteString(string(r))
		dashing = false
	}

	for _, r := range s {
		if r >= 'a' && r <= 'z' {
			out(r)
			continue
		}
		if r >= 'A' && r <= 'Z' {
			out(r - 'A' + 'a')
			continue
		}
		if r >= '0' && r <= '9' {
			out(r)
			continue
		}
		if buf.Len() > 0 {
			dashing = true
		}
	}

	ret := buf.String()
	if ret == "" {
		hash := sha256.Sum256([]byte(s))
		return hex.EncodeToString(hash[:])
	}
	return ret
}
