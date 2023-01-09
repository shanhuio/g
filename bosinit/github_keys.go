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

package bosinit

import (
	"fmt"
	"net/url"
	"strings"

	"shanhu.io/pub/httputil"
)

// FetchGitHubKeys fetches github ssh public keys of a user.
func FetchGitHubKeys(user string) ([]string, error) {
	c := &httputil.Client{
		Server: &url.URL{
			Scheme: "https",
			Host:   "github.com",
		},
	}

	keys, err := c.GetString(fmt.Sprintf("/%s.keys", user))
	if err != nil {
		return nil, err
	}

	var lines []string
	for _, line := range strings.Split(keys, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines, nil
}
