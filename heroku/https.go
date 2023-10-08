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

package heroku

import (
	"shanhu.io/g/aries"
)

// RedirectHTTPS redirects incoming HTTPS requests to HTTPS.
func RedirectHTTPS(c *aries.C) bool {
	if c.HTTPS {
		return false
	}

	u := c.Req.URL
	u.Host = c.Req.Host
	u.Scheme = "https"
	c.Redirect(u.String())
	return true
}
