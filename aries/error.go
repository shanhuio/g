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

package aries

import (
	"shanhu.io/pub/errcode"
)

const nothingHere = "nothing here"

// Miss is returned when a mux or router does not
// hit anything in its path lookup.
var Miss error = errcode.NotFoundf(nothingHere)

// NotFound is a true not found error.
var NotFound error = errcode.NotFoundf(nothingHere)

// NeedSignIn is returned when sign in is required for visiting a particular
// page.
var NeedSignIn error = errcode.Unauthorizedf("please sign in")
