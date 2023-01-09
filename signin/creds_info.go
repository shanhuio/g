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

package signin

import (
	"shanhu.io/pub/aries"
)

// CredsInfo is the user credential information got from gate checking.
type CredsInfo struct {
	Valid       bool
	NeedRefresh bool

	TokenType string
	User      string
	UserLevel int

	Data interface{}
}

const userDataKey = "user"

// ApplyCredsInfo applies the credential into the context.
func ApplyCredsInfo(c *aries.C, info *CredsInfo) {
	if !info.Valid {
		c.User = ""
		c.UserLevel = 0
		return
	}

	c.User = info.User
	c.UserLevel = info.UserLevel
	if info.Data != nil {
		c.Data[userDataKey] = info.Data
	}
}

// UserData fetches the user data in the context.
func UserData(c *aries.C) interface{} {
	v, ok := c.Data[userDataKey]
	if !ok {
		return nil
	}
	return v
}
