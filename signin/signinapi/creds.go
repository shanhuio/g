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

package signinapi

import (
	"time"

	"shanhu.io/pub/timeutil"
)

// Creds is the response for signing in. It saves the user credentials.
type Creds struct {
	User        string
	Token       string
	ExpiresTime *timeutil.Timestamp `json:",omitempty"`

	Expires int64 `json:",omitempty"` // Nanosecond timestamp, legacy use.
}

// FixTime fixes timestamps.
func (c *Creds) FixTime() {
	if c.ExpiresTime == nil && c.Expires != 0 {
		t := time.Unix(0, c.Expires)
		c.ExpiresTime = timeutil.NewTimestamp(t)
	}
}
