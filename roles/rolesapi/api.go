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

package rolesapi

import (
	"shanhu.io/pub/timeutil"
)

// Role contains the roles meta information.
type Role struct {
	Name       string
	TimeCreate *timeutil.Timestamp
	Disabled   bool `json:",omitempty"`
}

// PassCode contains the info of a pass code.
type PassCode struct {
	Code   string
	Valid  *timeutil.Timestamp
	Expire *timeutil.Timestamp

	TriedTooManyTimes bool
}