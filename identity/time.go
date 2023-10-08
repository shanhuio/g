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

package identity

import (
	"time"

	"shanhu.io/g/errcode"
)

func publicKeyValid(k *PublicKey, now time.Time) error {
	if k.NotValidBefore > 0 {
		if now.Before(time.Unix(k.NotValidBefore, 0)) {
			return errcode.InvalidArgf("key not valid yet")
		}
	}
	if now.After(time.Unix(k.NotValidAfter, 0)) {
		return errcode.InvalidArgf("key expired")
	}
	return nil
}
