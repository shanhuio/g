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

package backapi

import (
	"shanhu.io/pub/aries"
	"shanhu.io/pub/pisces"
)

// Router provides an API service router for the given PsqlTables.
func Router(b *pisces.Tables) *aries.Router {
	r := aries.NewRouter()

	r.File("create", func(*aries.C) error { return b.Create() })
	r.File("create-missing", func(*aries.C) error {
		return b.CreateMissing()
	})
	r.File("destroy", func(*aries.C) error { return b.Destroy() })

	return r
}
